package network

import (
	"github.com/sakaeshinya/tui-chat/internal/domain"
)

// NetworkPort はネットワーク通信のインターフェース
type NetworkPort interface {
	// SendMessage はメッセージを送信するメソッド
	SendMessage(msg domain.Message) error

	// ReceiveMessage はメッセージを受信するメソッド
	ReceiveMessage() (domain.Message, error)

	// Connect はサーバーに接続するメソッド
	Connect(address string) error

	// Listen はサーバーを起動するメソッド
	Listen(address string) error

	// Close は接続を閉じるメソッド
	Close() error
}

// UserInfo はユーザー情報の交換に使用する構造体
type UserInfo struct {
	ID   domain.UserID
	Name string
}

// ConnectionStatus は接続状態を表す型
type ConnectionStatus int

const (
	// StatusDisconnected は未接続状態
	StatusDisconnected ConnectionStatus = iota

	// StatusConnecting は接続中状態
	StatusConnecting

	// StatusConnected は接続済み状態
	StatusConnected

	// StatusError はエラー状態
	StatusError
)

// String はConnectionStatusの文字列表現を返すメソッド
func (s ConnectionStatus) String() string {
	switch s {
	case StatusDisconnected:
		return "未接続"
	case StatusConnecting:
		return "接続中"
	case StatusConnected:
		return "接続済み"
	case StatusError:
		return "エラー"
	default:
		return "不明"
	}
}

// MessageType はメッセージの種類を表す型
type MessageType int

const (
	// TypeChatMessage は通常のチャットメッセージ
	TypeChatMessage MessageType = iota

	// TypeUserInfo はユーザー情報メッセージ
	TypeUserInfo

	// TypeConnectionAck は接続確認メッセージ
	TypeConnectionAck

	// TypeDisconnect は切断メッセージ
	TypeDisconnect
)

// NetworkMessage はネットワーク経由で送受信するメッセージの構造体
type NetworkMessage struct {
	Type    MessageType
	Payload interface{}
}
