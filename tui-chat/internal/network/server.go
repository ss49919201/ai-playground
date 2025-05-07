package network

import (
	"encoding/gob"
	"errors"
	"net"
	"sync"

	"github.com/sakaeshinya/tui-chat/internal/domain"
)

// TCPServer はTCPベースのサーバー実装
type TCPServer struct {
	listener     net.Listener
	conn         net.Conn
	encoder      *gob.Encoder
	decoder      *gob.Decoder
	status       ConnectionStatus
	statusMutex  sync.RWMutex
	messageChan  chan domain.Message
	errorChan    chan error
	stopChan     chan struct{}
	running      bool
	runningMutex sync.RWMutex
}

// NewTCPServer はTCPサーバーを生成するファクトリ関数
func NewTCPServer() *TCPServer {
	return &TCPServer{
		status:      StatusDisconnected,
		messageChan: make(chan domain.Message, 100),
		errorChan:   make(chan error, 10),
		stopChan:    make(chan struct{}),
	}
}

// Listen はTCPサーバーを起動するメソッド
func (s *TCPServer) Listen(address string) error {
	s.setStatus(StatusConnecting)

	var err error
	s.listener, err = net.Listen("tcp", address)
	if err != nil {
		s.setStatus(StatusError)
		return err
	}

	s.setRunning(true)
	go s.acceptConnection()

	return nil
}

// acceptConnection はクライアント接続を受け付けるゴルーチン
func (s *TCPServer) acceptConnection() {
	defer s.setRunning(false)

	// クライアント接続を待機
	conn, err := s.listener.Accept()
	if err != nil {
		s.errorChan <- err
		s.setStatus(StatusError)
		return
	}

	s.conn = conn
	s.encoder = gob.NewEncoder(s.conn)
	s.decoder = gob.NewDecoder(s.conn)
	s.setStatus(StatusConnected)

	// 受信ループを開始
	go s.receiveLoop()
}

// receiveLoop はメッセージ受信ループ
func (s *TCPServer) receiveLoop() {
	defer func() {
		if s.conn != nil {
			s.conn.Close()
		}
		s.setStatus(StatusDisconnected)
	}()

	for {
		select {
		case <-s.stopChan:
			return
		default:
			var netMsg NetworkMessage
			err := s.decoder.Decode(&netMsg)
			if err != nil {
				s.errorChan <- err
				return
			}

			switch netMsg.Type {
			case TypeChatMessage:
				if msg, ok := netMsg.Payload.(domain.Message); ok {
					s.messageChan <- msg
				}
			case TypeDisconnect:
				return
			}
		}
	}
}

// SendMessage はメッセージを送信するメソッド
func (s *TCPServer) SendMessage(msg domain.Message) error {
	if s.getStatus() != StatusConnected {
		return errors.New("サーバーが接続されていません")
	}

	netMsg := NetworkMessage{
		Type:    TypeChatMessage,
		Payload: msg,
	}

	return s.encoder.Encode(netMsg)
}

// ReceiveMessage はメッセージを受信するメソッド
func (s *TCPServer) ReceiveMessage() (domain.Message, error) {
	select {
	case msg := <-s.messageChan:
		return msg, nil
	case err := <-s.errorChan:
		return domain.Message{}, err
	}
}

// Connect はサーバーモードでは使用しないメソッド
func (s *TCPServer) Connect(address string) error {
	return errors.New("サーバーモードではConnectメソッドは使用できません")
}

// Close はサーバーを終了するメソッド
func (s *TCPServer) Close() error {
	if s.isRunning() {
		close(s.stopChan)
	}

	if s.conn != nil {
		s.conn.Close()
	}

	if s.listener != nil {
		s.listener.Close()
	}

	s.setStatus(StatusDisconnected)
	return nil
}

// GetStatus は現在の接続状態を取得するメソッド
func (s *TCPServer) GetStatus() ConnectionStatus {
	return s.getStatus()
}

// setStatus は接続状態を設定する内部メソッド
func (s *TCPServer) setStatus(status ConnectionStatus) {
	s.statusMutex.Lock()
	defer s.statusMutex.Unlock()
	s.status = status
}

// getStatus は接続状態を取得する内部メソッド
func (s *TCPServer) getStatus() ConnectionStatus {
	s.statusMutex.RLock()
	defer s.statusMutex.RUnlock()
	return s.status
}

// setRunning は実行状態を設定する内部メソッド
func (s *TCPServer) setRunning(running bool) {
	s.runningMutex.Lock()
	defer s.runningMutex.Unlock()
	s.running = running
}

// isRunning は実行状態を取得する内部メソッド
func (s *TCPServer) isRunning() bool {
	s.runningMutex.RLock()
	defer s.runningMutex.RUnlock()
	return s.running
}

func init() {
	// gob登録
	gob.Register(domain.Message{})
	gob.Register(domain.UserID(""))
	gob.Register(UserInfo{})
	gob.Register(NetworkMessage{})
}
