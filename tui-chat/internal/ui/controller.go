package ui

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sakaeshinya/tui-chat/internal/domain"
)

// TUIController はTUIコントローラ
type TUIController struct {
	app        *tview.Application
	chatView   *tview.TextView
	inputField *tview.InputField
	statusBar  *tview.TextView
	rootFlex   *tview.Flex
	handler    UIEventHandler
	user       domain.User
	mode       UIMode
	modeMutex  sync.RWMutex
	stopChan   chan struct{}
}

// NewTUIController はTUIコントローラを生成するファクトリ関数
func NewTUIController(handler UIEventHandler, user domain.User) *TUIController {
	controller := &TUIController{
		app:      tview.NewApplication(),
		handler:  handler,
		user:     user,
		mode:     ModeNormal,
		stopChan: make(chan struct{}),
	}

	controller.setupUI()
	return controller
}

// setupUI はUIコンポーネントを設定するメソッド
func (c *TUIController) setupUI() {
	// チャット表示領域
	c.chatView = tview.NewTextView().
		SetDynamicColors(true).
		SetChangedFunc(func() {
			c.app.Draw()
		})
	c.chatView.SetBorder(true).SetTitle("チャット")
	c.chatView.SetScrollable(true)
	c.chatView.SetWordWrap(true)

	// 入力領域
	c.inputField = tview.NewInputField().
		SetLabel("> ").
		SetFieldWidth(0).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				text := c.inputField.GetText()
				if text != "" {
					c.handler.OnMessageSend(text)
					c.inputField.SetText("")
				}
			}
		})
	c.inputField.SetBorder(true).SetTitle("メッセージ入力")

	// ステータスバー
	c.statusBar = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	c.statusBar.SetBorder(true).SetTitle("ステータス")

	// コマンドヘルプ
	helpText := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	helpText.SetText("[yellow]Ctrl+S[white]: サーバー開始 | [yellow]Ctrl+C[white]: クライアント接続 | [yellow]Ctrl+D[white]: 切断 | [yellow]Ctrl+Q[white]: 終了")
	helpText.SetBorder(false)

	// レイアウト設定
	c.rootFlex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(c.chatView, 0, 1, false).
		AddItem(c.inputField, 3, 0, true).
		AddItem(c.statusBar, 3, 0, false).
		AddItem(helpText, 1, 0, false)

	// キーバインディング
	c.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlS:
			// サーバーモード
			c.showServerDialog()
			return nil
		case tcell.KeyCtrlC:
			// クライアントモード
			c.showClientDialog()
			return nil
		case tcell.KeyCtrlD:
			// 切断
			c.handler.OnDisconnect()
			return nil
		case tcell.KeyCtrlQ:
			// 終了
			c.handler.OnQuit()
			return nil
		}
		return event
	})

	c.app.SetRoot(c.rootFlex, true)
}

// showServerDialog はサーバーモードのダイアログを表示するメソッド
func (c *TUIController) showServerDialog() {
	form := tview.NewForm()
	form.AddInputField("ポート番号", "8080", 20, nil, nil)
	form.AddButton("開始", func() {
		port := form.GetFormItem(0).(*tview.InputField).GetText()
		if port == "" {
			port = "8080"
		}
		address := fmt.Sprintf(":%s", port)
		c.handler.OnListen(address)
		c.app.SetRoot(c.rootFlex, true)
	})
	form.AddButton("キャンセル", func() {
		c.app.SetRoot(c.rootFlex, true)
	})
	form.SetBorder(true).SetTitle("サーバーモード")
	form.SetCancelFunc(func() {
		c.app.SetRoot(c.rootFlex, true)
	})

	c.app.SetRoot(form, true)
}

// showClientDialog はクライアントモードのダイアログを表示するメソッド
func (c *TUIController) showClientDialog() {
	form := tview.NewForm()
	form.AddInputField("サーバーアドレス", "localhost:8080", 30, nil, nil)
	form.AddButton("接続", func() {
		address := form.GetFormItem(0).(*tview.InputField).GetText()
		if address == "" {
			address = "localhost:8080"
		}
		if !strings.Contains(address, ":") {
			address = address + ":8080"
		}
		c.handler.OnConnect(address)
		c.app.SetRoot(c.rootFlex, true)
	})
	form.AddButton("キャンセル", func() {
		c.app.SetRoot(c.rootFlex, true)
	})
	form.SetBorder(true).SetTitle("クライアントモード")
	form.SetCancelFunc(func() {
		c.app.SetRoot(c.rootFlex, true)
	})

	c.app.SetRoot(form, true)
}

// DisplayMessage はメッセージを表示するメソッド
func (c *TUIController) DisplayMessage(msg domain.Message) {
	c.chatView.Write([]byte(fmt.Sprintf("[%s] %s: %s\n",
		msg.FormattedTime(),
		msg.Sender,
		msg.Content)))
}

// GetInput は入力を取得するメソッド
func (c *TUIController) GetInput() (string, error) {
	return c.inputField.GetText(), nil
}

// UpdateStatus はステータスを更新するメソッド
func (c *TUIController) UpdateStatus(status string) {
	c.statusBar.Clear()
	c.statusBar.Write([]byte(status))
}

// ShowError はエラーを表示するメソッド
func (c *TUIController) ShowError(err error) {
	c.statusBar.Clear()
	c.statusBar.Write([]byte(fmt.Sprintf("[red]エラー: %s[-]", err.Error())))
}

// Start はTUIを開始するメソッド
func (c *TUIController) Start() error {
	// 初期メッセージ
	c.chatView.Write([]byte("[yellow]TUIチャットアプリケーションへようこそ！[-]\n"))
	c.chatView.Write([]byte("[yellow]Ctrl+S[white]でサーバーを起動するか、[yellow]Ctrl+C[white]でサーバーに接続してください。\n"))
	c.UpdateStatus(fmt.Sprintf("ユーザー: %s | モード: %s", c.user.Name, c.getMode().String()))

	return c.app.Run()
}

// Stop はTUIを停止するメソッド
func (c *TUIController) Stop() {
	c.app.Stop()
}

// SetMode はUIモードを設定するメソッド
func (c *TUIController) SetMode(mode UIMode) {
	c.modeMutex.Lock()
	defer c.modeMutex.Unlock()
	c.mode = mode
	c.UpdateStatus(fmt.Sprintf("ユーザー: %s | モード: %s", c.user.Name, mode.String()))
}

// getMode はUIモードを取得する内部メソッド
func (c *TUIController) getMode() UIMode {
	c.modeMutex.RLock()
	defer c.modeMutex.RUnlock()
	return c.mode
}

// SetHandler はイベントハンドラを設定するメソッド
func (c *TUIController) SetHandler(handler UIEventHandler) {
	c.handler = handler
}
