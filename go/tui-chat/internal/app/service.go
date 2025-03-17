package app

import (
	"fmt"
	"sync"

	"github.com/sakaeshinya/tui-chat/internal/domain"
	"github.com/sakaeshinya/tui-chat/internal/network"
	"github.com/sakaeshinya/tui-chat/internal/ui"
)

// ChatService はチャットアプリケーションのコアサービス
type ChatService struct {
	session      *domain.ChatSession
	sessionMutex sync.RWMutex
	messageQueue chan domain.Message
	network      network.NetworkPort
	ui           ui.UIPort
	user         domain.User
	stopChan     chan struct{}
	running      bool
	runningMutex sync.RWMutex
}

// NewChatService はチャットサービスを生成するファクトリ関数
func NewChatService(network network.NetworkPort, ui ui.UIPort, user domain.User) *ChatService {
	return &ChatService{
		messageQueue: make(chan domain.Message, 100),
		network:      network,
		ui:           ui,
		user:         user,
		stopChan:     make(chan struct{}),
	}
}

// StartServer はサーバーモードでチャットを開始するメソッド
func (s *ChatService) StartServer(address string) error {
	// セッション初期化
	s.sessionMutex.Lock()
	s.session = new(domain.ChatSession)
	*s.session = domain.NewChatSession(s.user)
	s.sessionMutex.Unlock()

	// サーバー起動
	err := s.network.Listen(address)
	if err != nil {
		return err
	}

	s.setRunning(true)
	go s.processIncomingMessages()
	go s.processOutgoingMessages()

	s.ui.UpdateStatus(fmt.Sprintf("サーバーモード | アドレス: %s", address))
	return nil
}

// JoinChat はクライアントモードでチャットに参加するメソッド
func (s *ChatService) JoinChat(address string) error {
	// クライアント接続
	err := s.network.Connect(address)
	if err != nil {
		return err
	}

	// セッション初期化（仮）
	s.sessionMutex.Lock()
	s.session = new(domain.ChatSession)
	*s.session = domain.NewChatSession(s.user)
	s.sessionMutex.Unlock()

	s.setRunning(true)
	go s.processIncomingMessages()
	go s.processOutgoingMessages()

	s.ui.UpdateStatus(fmt.Sprintf("クライアントモード | 接続先: %s", address))
	return nil
}

// SendMessage はメッセージを送信するメソッド
func (s *ChatService) SendMessage(content string) error {
	msg, err := domain.NewMessage(content, s.user.ID)
	if err != nil {
		return err
	}

	s.messageQueue <- msg
	return nil
}

// Disconnect は接続を切断するメソッド
func (s *ChatService) Disconnect() error {
	if !s.isRunning() {
		return nil
	}

	close(s.stopChan)
	s.setRunning(false)

	err := s.network.Close()
	if err != nil {
		return err
	}

	s.ui.UpdateStatus("切断しました")
	return nil
}

// processIncomingMessages は受信メッセージを処理するメソッド
func (s *ChatService) processIncomingMessages() {
	defer s.setRunning(false)

	for {
		select {
		case <-s.stopChan:
			return
		default:
			msg, err := s.network.ReceiveMessage()
			if err != nil {
				s.ui.ShowError(err)
				return
			}

			s.sessionMutex.Lock()
			err = s.session.AddMessage(msg)
			s.sessionMutex.Unlock()

			if err != nil {
				s.ui.ShowError(err)
				continue
			}

			s.ui.DisplayMessage(msg)
		}
	}
}

// processOutgoingMessages は送信メッセージを処理するメソッド
func (s *ChatService) processOutgoingMessages() {
	for {
		select {
		case <-s.stopChan:
			return
		case msg := <-s.messageQueue:
			err := s.network.SendMessage(msg)
			if err != nil {
				s.ui.ShowError(err)
				continue
			}

			s.sessionMutex.Lock()
			err = s.session.AddMessage(msg)
			s.sessionMutex.Unlock()

			if err != nil {
				s.ui.ShowError(err)
				continue
			}

			s.ui.DisplayMessage(msg)
		}
	}
}

// setRunning は実行状態を設定する内部メソッド
func (s *ChatService) setRunning(running bool) {
	s.runningMutex.Lock()
	defer s.runningMutex.Unlock()
	s.running = running
}

// isRunning は実行状態を取得する内部メソッド
func (s *ChatService) isRunning() bool {
	s.runningMutex.RLock()
	defer s.runningMutex.RUnlock()
	return s.running
}
