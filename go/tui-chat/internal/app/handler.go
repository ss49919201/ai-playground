package app

import (
	"os"

	"github.com/sakaeshinya/tui-chat/internal/ui"
)

// ChatEventHandler はチャットイベントを処理するハンドラ
type ChatEventHandler struct {
	service *ChatService
	ui      ui.UIPort
}

// NewChatEventHandler はチャットイベントハンドラを生成するファクトリ関数
func NewChatEventHandler(service *ChatService, ui ui.UIPort) *ChatEventHandler {
	return &ChatEventHandler{
		service: service,
		ui:      ui,
	}
}

// OnMessageSend はメッセージ送信イベントを処理するメソッド
func (h *ChatEventHandler) OnMessageSend(content string) {
	err := h.service.SendMessage(content)
	if err != nil {
		h.ui.ShowError(err)
	}
}

// OnConnect は接続イベントを処理するメソッド
func (h *ChatEventHandler) OnConnect(address string) {
	err := h.service.JoinChat(address)
	if err != nil {
		h.ui.ShowError(err)
		return
	}

	// UIモードを更新
	if controller, ok := h.ui.(*ui.TUIController); ok {
		controller.SetMode(ui.ModeClient)
	}
}

// OnListen はリッスンイベントを処理するメソッド
func (h *ChatEventHandler) OnListen(address string) {
	err := h.service.StartServer(address)
	if err != nil {
		h.ui.ShowError(err)
		return
	}

	// UIモードを更新
	if controller, ok := h.ui.(*ui.TUIController); ok {
		controller.SetMode(ui.ModeServer)
	}
}

// OnDisconnect は切断イベントを処理するメソッド
func (h *ChatEventHandler) OnDisconnect() {
	err := h.service.Disconnect()
	if err != nil {
		h.ui.ShowError(err)
		return
	}

	// UIモードを更新
	if controller, ok := h.ui.(*ui.TUIController); ok {
		controller.SetMode(ui.ModeNormal)
	}
}

// OnQuit は終了イベントを処理するメソッド
func (h *ChatEventHandler) OnQuit() {
	// 切断処理
	h.service.Disconnect()

	// UIを停止
	h.ui.Stop()

	// アプリケーション終了
	os.Exit(0)
}
