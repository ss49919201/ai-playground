package ui

import (
	"github.com/sakaeshinya/tui-chat/internal/domain"
)

// UIPort はユーザーインターフェースのインターフェース
type UIPort interface {
	// DisplayMessage はメッセージを表示するメソッド
	DisplayMessage(msg domain.Message)

	// GetInput は入力を取得するメソッド
	GetInput() (string, error)

	// UpdateStatus はステータスを更新するメソッド
	UpdateStatus(status string)

	// ShowError はエラーを表示するメソッド
	ShowError(err error)

	// Start はUIを開始するメソッド
	Start() error

	// Stop はUIを停止するメソッド
	Stop()
}

// UIEventHandler はUIイベントを処理するハンドラのインターフェース
type UIEventHandler interface {
	// OnMessageSend はメッセージ送信イベントを処理するメソッド
	OnMessageSend(content string)

	// OnConnect は接続イベントを処理するメソッド
	OnConnect(address string)

	// OnListen はリッスンイベントを処理するメソッド
	OnListen(address string)

	// OnDisconnect は切断イベントを処理するメソッド
	OnDisconnect()

	// OnQuit は終了イベントを処理するメソッド
	OnQuit()
}

// UIMode はUIのモードを表す型
type UIMode int

const (
	// ModeNormal は通常モード
	ModeNormal UIMode = iota

	// ModeServer はサーバーモード
	ModeServer

	// ModeClient はクライアントモード
	ModeClient
)

// String はUIModeの文字列表現を返すメソッド
func (m UIMode) String() string {
	switch m {
	case ModeNormal:
		return "通常モード"
	case ModeServer:
		return "サーバーモード"
	case ModeClient:
		return "クライアントモード"
	default:
		return "不明なモード"
	}
}
