package network

import (
	"encoding/gob"
	"errors"
	"net"
	"sync"

	"github.com/sakaeshinya/tui-chat/internal/domain"
)

// TCPClient はTCPベースのクライアント実装
type TCPClient struct {
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

// NewTCPClient はTCPクライアントを生成するファクトリ関数
func NewTCPClient() *TCPClient {
	return &TCPClient{
		status:      StatusDisconnected,
		messageChan: make(chan domain.Message, 100),
		errorChan:   make(chan error, 10),
		stopChan:    make(chan struct{}),
	}
}

// Connect はサーバーに接続するメソッド
func (c *TCPClient) Connect(address string) error {
	c.setStatus(StatusConnecting)

	var err error
	c.conn, err = net.Dial("tcp", address)
	if err != nil {
		c.setStatus(StatusError)
		return err
	}

	c.encoder = gob.NewEncoder(c.conn)
	c.decoder = gob.NewDecoder(c.conn)
	c.setStatus(StatusConnected)

	c.setRunning(true)
	go c.receiveLoop()

	return nil
}

// receiveLoop はメッセージ受信ループ
func (c *TCPClient) receiveLoop() {
	defer func() {
		if c.conn != nil {
			c.conn.Close()
		}
		c.setStatus(StatusDisconnected)
		c.setRunning(false)
	}()

	for {
		select {
		case <-c.stopChan:
			return
		default:
			var netMsg NetworkMessage
			err := c.decoder.Decode(&netMsg)
			if err != nil {
				c.errorChan <- err
				return
			}

			switch netMsg.Type {
			case TypeChatMessage:
				if msg, ok := netMsg.Payload.(domain.Message); ok {
					c.messageChan <- msg
				}
			case TypeDisconnect:
				return
			}
		}
	}
}

// SendMessage はメッセージを送信するメソッド
func (c *TCPClient) SendMessage(msg domain.Message) error {
	if c.getStatus() != StatusConnected {
		return errors.New("クライアントが接続されていません")
	}

	netMsg := NetworkMessage{
		Type:    TypeChatMessage,
		Payload: msg,
	}

	return c.encoder.Encode(netMsg)
}

// ReceiveMessage はメッセージを受信するメソッド
func (c *TCPClient) ReceiveMessage() (domain.Message, error) {
	select {
	case msg := <-c.messageChan:
		return msg, nil
	case err := <-c.errorChan:
		return domain.Message{}, err
	}
}

// Listen はクライアントモードでは使用しないメソッド
func (c *TCPClient) Listen(address string) error {
	return errors.New("クライアントモードではListenメソッドは使用できません")
}

// Close はクライアントを終了するメソッド
func (c *TCPClient) Close() error {
	if c.isRunning() {
		// 切断メッセージを送信
		netMsg := NetworkMessage{
			Type: TypeDisconnect,
		}
		c.encoder.Encode(netMsg)

		close(c.stopChan)
	}

	if c.conn != nil {
		c.conn.Close()
	}

	c.setStatus(StatusDisconnected)
	return nil
}

// GetStatus は現在の接続状態を取得するメソッド
func (c *TCPClient) GetStatus() ConnectionStatus {
	return c.getStatus()
}

// SendUserInfo はユーザー情報を送信するメソッド
func (c *TCPClient) SendUserInfo(user domain.User) error {
	if c.getStatus() != StatusConnected {
		return errors.New("クライアントが接続されていません")
	}

	userInfo := UserInfo{
		ID:   user.ID,
		Name: user.Name,
	}

	netMsg := NetworkMessage{
		Type:    TypeUserInfo,
		Payload: userInfo,
	}

	return c.encoder.Encode(netMsg)
}

// setStatus は接続状態を設定する内部メソッド
func (c *TCPClient) setStatus(status ConnectionStatus) {
	c.statusMutex.Lock()
	defer c.statusMutex.Unlock()
	c.status = status
}

// getStatus は接続状態を取得する内部メソッド
func (c *TCPClient) getStatus() ConnectionStatus {
	c.statusMutex.RLock()
	defer c.statusMutex.RUnlock()
	return c.status
}

// setRunning は実行状態を設定する内部メソッド
func (c *TCPClient) setRunning(running bool) {
	c.runningMutex.Lock()
	defer c.runningMutex.Unlock()
	c.running = running
}

// isRunning は実行状態を取得する内部メソッド
func (c *TCPClient) isRunning() bool {
	c.runningMutex.RLock()
	defer c.runningMutex.RUnlock()
	return c.running
}
