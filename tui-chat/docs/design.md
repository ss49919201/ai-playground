# TUI チャットアプリケーション設計書

## 1. アプリケーション概要

このアプリケーションは、テキストベースのユーザーインターフェース（TUI）を使用したチャットアプリケーションです。TCP 通信を使用して、最大 2 人のユーザー間でのリアルタイムなメッセージのやり取りを可能にします。サーバーモードとクライアントモードの両方を提供し、ローカルネットワーク内での通信をサポートします。

## 2. アーキテクチャ設計

アプリケーションは以下の主要なコンポーネントで構成されます：

```
+----------------+      +----------------+      +----------------+
|                |      |                |      |                |
|  UI コンポーネント  <----->  アプリケーション  <----->  通信コンポーネント |
|                |      |    コンポーネント   |      |                |
+----------------+      +----------------+      +----------------+
                                ^
                                |
                                v
                        +----------------+
                        |                |
                        |  モデルコンポーネント |
                        |                |
                        +----------------+
```

- **UI コンポーネント**: TUI の表示と入力処理を担当
- **アプリケーションコンポーネント**: ビジネスロジックとフロー制御を担当
- **通信コンポーネント**: TCP 通信の確立と維持を担当
- **モデルコンポーネント**: データモデルとドメインロジックを担当

### 2.1 アーキテクチャパターン

- ヘキサゴナルアーキテクチャ（ポートとアダプターパターン）を採用
- 関数型アプローチを重視し、純粋関数と副作用の分離を実現
- 依存性の注入を活用して、コンポーネント間の結合度を低減

## 3. モジュール構成

アプリケーションは以下のパッケージ構造で実装します：

```
tui-chat/
├── cmd/
│   └── chat/
│       └── main.go           # エントリーポイント
├── internal/
│   ├── domain/               # ドメインモデルとビジネスロジック
│   │   ├── message.go        # メッセージモデル
│   │   ├── user.go           # ユーザーモデル
│   │   └── chat.go           # チャットセッションモデル
│   ├── app/                  # アプリケーションサービス
│   │   ├── service.go        # アプリケーションサービス
│   │   └── handler.go        # イベントハンドラ
│   ├── ui/                   # ユーザーインターフェース
│   │   ├── tui.go            # TUIコンポーネント
│   │   ├── view.go           # ビュー定義
│   │   └── controller.go     # UIコントローラ
│   ├── network/              # ネットワーク通信
│   │   ├── server.go         # サーバー実装
│   │   ├── client.go         # クライアント実装
│   │   └── protocol.go       # 通信プロトコル定義
│   └── config/               # 設定管理
│       └── config.go         # 設定定義
└── pkg/                      # 再利用可能なパッケージ
    └── utils/                # ユーティリティ関数
        └── logger.go         # ロギングユーティリティ
```

## 4. クラス設計

### 4.1 ドメインモデル

#### 4.1.1 Message（メッセージ）

```go
// Message はチャットメッセージを表す値オブジェクト
type Message struct {
    ID        string    // メッセージの一意識別子
    Content   string    // メッセージ内容
    Sender    UserID    // 送信者ID
    Timestamp time.Time // 送信時刻
}

// NewMessage はメッセージオブジェクトを生成するファクトリ関数
func NewMessage(content string, sender UserID) (Message, error) {
    if content == "" {
        return Message{}, errors.New("メッセージ内容が空です")
    }
    return Message{
        ID:        uuid.New().String(),
        Content:   content,
        Sender:    sender,
        Timestamp: time.Now(),
    }, nil
}
```

#### 4.1.2 User（ユーザー）

```go
// UserID はユーザーIDを表すカスタム型
type UserID string

// User はチャットユーザーを表すエンティティ
type User struct {
    ID       UserID
    Name     string
    IsActive bool
}

// NewUser はユーザーオブジェクトを生成するファクトリ関数
func NewUser(name string) (User, error) {
    if name == "" {
        return User{}, errors.New("ユーザー名が空です")
    }
    return User{
        ID:       UserID(uuid.New().String()),
        Name:     name,
        IsActive: true,
    }, nil
}
```

#### 4.1.3 ChatSession（チャットセッション）

```go
// ChatSession はチャットセッションを表す集約ルート
type ChatSession struct {
    ID        string
    Users     []User
    Messages  []Message
    CreatedAt time.Time
    UpdatedAt time.Time
}

// NewChatSession はチャットセッションを生成するファクトリ関数
func NewChatSession(initialUser User) ChatSession {
    return ChatSession{
        ID:        uuid.New().String(),
        Users:     []User{initialUser},
        Messages:  []Message{},
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
}

// AddUser はチャットセッションにユーザーを追加するメソッド
func (c *ChatSession) AddUser(user User) error {
    if len(c.Users) >= 2 {
        return errors.New("チャットセッションは最大2人までです")
    }
    c.Users = append(c.Users, user)
    c.UpdatedAt = time.Now()
    return nil
}

// AddMessage はチャットセッションにメッセージを追加するメソッド
func (c *ChatSession) AddMessage(msg Message) error {
    // 送信者がセッションに参加しているか確認
    senderExists := false
    for _, user := range c.Users {
        if user.ID == msg.Sender {
            senderExists = true
            break
        }
    }
    if !senderExists {
        return errors.New("送信者がチャットセッションに参加していません")
    }

    c.Messages = append(c.Messages, msg)
    c.UpdatedAt = time.Now()
    return nil
}
```

### 4.2 アプリケーションサービス

#### 4.2.1 ChatService（チャットサービス）

```go
// ChatService はチャットアプリケーションのコアサービス
type ChatService struct {
    session      *ChatSession
    messageQueue chan Message
    network      NetworkPort
    ui           UIPort
}

// NetworkPort はネットワーク通信のインターフェース
type NetworkPort interface {
    SendMessage(msg Message) error
    ReceiveMessage() (Message, error)
    Connect(address string) error
    Listen(address string) error
    Close() error
}

// UIPort はユーザーインターフェースのインターフェース
type UIPort interface {
    DisplayMessage(msg Message)
    GetInput() (string, error)
    UpdateStatus(status string)
    ShowError(err error)
}

// NewChatService はチャットサービスを生成するファクトリ関数
func NewChatService(network NetworkPort, ui UIPort) *ChatService {
    return &ChatService{
        session:      nil,
        messageQueue: make(chan Message, 100),
        network:      network,
        ui:           ui,
    }
}

// StartServer はサーバーモードでチャットを開始するメソッド
func (s *ChatService) StartServer(user User, address string) error {
    s.session = &ChatSession{}
    *s.session = NewChatSession(user)

    err := s.network.Listen(address)
    if err != nil {
        return err
    }

    go s.processIncomingMessages()
    go s.processOutgoingMessages()

    return nil
}

// JoinChat はクライアントモードでチャットに参加するメソッド
func (s *ChatService) JoinChat(user User, address string) error {
    err := s.network.Connect(address)
    if err != nil {
        return err
    }

    // 接続後にユーザー情報を送信
    // 実装省略

    go s.processIncomingMessages()
    go s.processOutgoingMessages()

    return nil
}

// SendMessage はメッセージを送信するメソッド
func (s *ChatService) SendMessage(content string, sender UserID) error {
    msg, err := NewMessage(content, sender)
    if err != nil {
        return err
    }

    s.messageQueue <- msg
    return nil
}

// processIncomingMessages は受信メッセージを処理するメソッド
func (s *ChatService) processIncomingMessages() {
    for {
        msg, err := s.network.ReceiveMessage()
        if err != nil {
            s.ui.ShowError(err)
            continue
        }

        err = s.session.AddMessage(msg)
        if err != nil {
            s.ui.ShowError(err)
            continue
        }

        s.ui.DisplayMessage(msg)
    }
}

// processOutgoingMessages は送信メッセージを処理するメソッド
func (s *ChatService) processOutgoingMessages() {
    for msg := range s.messageQueue {
        err := s.network.SendMessage(msg)
        if err != nil {
            s.ui.ShowError(err)
            continue
        }

        err = s.session.AddMessage(msg)
        if err != nil {
            s.ui.ShowError(err)
            continue
        }

        s.ui.DisplayMessage(msg)
    }
}
```

### 4.3 ネットワークコンポーネント

#### 4.3.1 TCPServer（TCP サーバー）

```go
// TCPServer はTCPベースのサーバー実装
type TCPServer struct {
    listener net.Listener
    conn     net.Conn
    encoder  *gob.Encoder
    decoder  *gob.Decoder
}

// NewTCPServer はTCPサーバーを生成するファクトリ関数
func NewTCPServer() *TCPServer {
    return &TCPServer{}
}

// Listen はTCPサーバーを起動するメソッド
func (s *TCPServer) Listen(address string) error {
    var err error
    s.listener, err = net.Listen("tcp", address)
    if err != nil {
        return err
    }

    // クライアント接続を待機
    s.conn, err = s.listener.Accept()
    if err != nil {
        return err
    }

    s.encoder = gob.NewEncoder(s.conn)
    s.decoder = gob.NewDecoder(s.conn)

    return nil
}

// SendMessage はメッセージを送信するメソッド
func (s *TCPServer) SendMessage(msg Message) error {
    return s.encoder.Encode(msg)
}

// ReceiveMessage はメッセージを受信するメソッド
func (s *TCPServer) ReceiveMessage() (Message, error) {
    var msg Message
    err := s.decoder.Decode(&msg)
    return msg, err
}

// Close はサーバーを終了するメソッド
func (s *TCPServer) Close() error {
    if s.conn != nil {
        s.conn.Close()
    }
    if s.listener != nil {
        s.listener.Close()
    }
    return nil
}
```

#### 4.3.2 TCPClient（TCP クライアント）

```go
// TCPClient はTCPベースのクライアント実装
type TCPClient struct {
    conn    net.Conn
    encoder *gob.Encoder
    decoder *gob.Decoder
}

// NewTCPClient はTCPクライアントを生成するファクトリ関数
func NewTCPClient() *TCPClient {
    return &TCPClient{}
}

// Connect はサーバーに接続するメソッド
func (c *TCPClient) Connect(address string) error {
    var err error
    c.conn, err = net.Dial("tcp", address)
    if err != nil {
        return err
    }

    c.encoder = gob.NewEncoder(c.conn)
    c.decoder = gob.NewDecoder(c.conn)

    return nil
}

// SendMessage はメッセージを送信するメソッド
func (c *TCPClient) SendMessage(msg Message) error {
    return c.encoder.Encode(msg)
}

// ReceiveMessage はメッセージを受信するメソッド
func (c *TCPClient) ReceiveMessage() (Message, error) {
    var msg Message
    err := c.decoder.Decode(&msg)
    return msg, err
}

// Close はクライアントを終了するメソッド
func (c *TCPClient) Close() error {
    if c.conn != nil {
        c.conn.Close()
    }
    return nil
}
```

### 4.4 UI コンポーネント

#### 4.4.1 TUIController（TUI コントローラ）

```go
// TUIController はTUIコントローラ
type TUIController struct {
    app        *tview.Application
    chatView   *tview.TextView
    inputField *tview.InputField
    statusBar  *tview.TextView
    service    *ChatService
    user       User
}

// NewTUIController はTUIコントローラを生成するファクトリ関数
func NewTUIController(service *ChatService, user User) *TUIController {
    controller := &TUIController{
        app:     tview.NewApplication(),
        service: service,
        user:    user,
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

    // 入力領域
    c.inputField = tview.NewInputField().
        SetLabel("> ").
        SetFieldWidth(0).
        SetDoneFunc(func(key tcell.Key) {
            if key == tcell.KeyEnter {
                text := c.inputField.GetText()
                if text != "" {
                    c.service.SendMessage(text, c.user.ID)
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

    // レイアウト設定
    flex := tview.NewFlex().
        SetDirection(tview.FlexRow).
        AddItem(c.chatView, 0, 1, false).
        AddItem(c.inputField, 3, 0, true).
        AddItem(c.statusBar, 3, 0, false)

    c.app.SetRoot(flex, true)
}

// Start はTUIを開始するメソッド
func (c *TUIController) Start() error {
    return c.app.Run()
}

// DisplayMessage はメッセージを表示するメソッド
func (c *TUIController) DisplayMessage(msg Message) {
    c.chatView.Write([]byte(fmt.Sprintf("[%s] %s: %s\n",
        msg.Timestamp.Format("15:04:05"),
        msg.Sender,
        msg.Content)))
}

// GetInput は入力を取得するメソッド
func (c *TUIController) GetInput() (string, error) {
    // 実際の実装では、入力はイベント駆動で処理されるため、
    // このメソッドは使用されない可能性があります
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
```

## 5. データフロー

### 5.1 メッセージ送信フロー

1. ユーザーが入力フィールドにテキストを入力し、Enter キーを押す
2. TUIController がテキストを取得し、ChatService の SendMessage メソッドを呼び出す
3. ChatService が Message オブジェクトを生成し、messageQueue に追加
4. processOutgoingMessages ゴルーチンがキューからメッセージを取り出し、NetworkPort の SendMessage メソッドを呼び出す
5. TCPServer または TCPClient がメッセージをエンコードし、ネットワーク経由で送信
6. ChatService がセッションにメッセージを追加し、UIPort の DisplayMessage メソッドを呼び出す
7. TUIController がメッセージをチャット表示領域に表示

### 5.2 メッセージ受信フロー

1. TCPServer または TCPClient がネットワーク経由でメッセージを受信
2. NetworkPort の ReceiveMessage メソッドがメッセージをデコードし、Message オブジェクトを返す
3. processIncomingMessages ゴルーチンがメッセージを受け取り、セッションに追加
4. UIPort の DisplayMessage メソッドを呼び出し、メッセージを表示
5. TUIController がメッセージをチャット表示領域に表示

## 6. ユーザーインターフェース設計

TUI は以下の 3 つの主要な領域で構成されます：

```
+---------------------------------------------------------------+
|                          チャット                              |
|                                                               |
| [14:30:45] ユーザー1: こんにちは                               |
| [14:31:02] ユーザー2: やあ、元気？                             |
| [14:31:15] ユーザー1: うん、元気だよ！                         |
|                                                               |
|                                                               |
+---------------------------------------------------------------+
|                        メッセージ入力                          |
| > _                                                           |
+---------------------------------------------------------------+
|                         ステータス                             |
| 接続済み: ユーザー1, ユーザー2                                 |
+---------------------------------------------------------------+
```

- **チャット領域**: 過去のメッセージを時系列で表示
- **メッセージ入力領域**: ユーザーが新しいメッセージを入力
- **ステータス領域**: 接続状態、エラーメッセージなどを表示

## 7. 通信プロトコル設計

通信プロトコルは TCP ベースで、メッセージの送受信には Go 言語の`gob`パッケージを使用します。これにより、Go の型をそのままシリアライズ/デシリアライズできます。

### 7.1 メッセージフォーマット

```go
// Message構造体がそのままシリアライズされます
type Message struct {
    ID        string    // メッセージの一意識別子
    Content   string    // メッセージ内容
    Sender    UserID    // 送信者ID
    Timestamp time.Time // 送信時刻
}
```

### 7.2 接続確立プロトコル

1. サーバーが特定のポートでリッスン開始
2. クライアントがサーバーに接続
3. クライアントがユーザー情報を送信
4. サーバーがユーザー情報を受信し、セッションにユーザーを追加
5. サーバーが接続確認メッセージを送信
6. 通常のメッセージ交換が開始

## 8. エラーハンドリング

エラーハンドリングは以下の原則に従います：

1. 各レイヤーで適切なエラーを生成し、上位レイヤーに伝播
2. ドメインエラーは明確なエラーメッセージを持つ専用の型で表現
3. ネットワークエラーは再接続ロジックで対応
4. UI レイヤーでユーザーにわかりやすいエラーメッセージを表示

```go
// エラー定義の例
var (
    ErrInvalidMessage   = errors.New("無効なメッセージです")
    ErrSessionFull      = errors.New("チャットセッションは満員です")
    ErrUserNotInSession = errors.New("ユーザーはセッションに参加していません")
    ErrNetworkFailure   = errors.New("ネットワーク接続に失敗しました")
)
```

## 9. テスト計画

テストは以下のレベルで実施します：

### 9.1 ユニットテスト

- ドメインモデルの各メソッドのテスト
- アプリケーションサービスのロジックテスト
- ネットワークコンポーネントのモックを使用したテスト

### 9.2 統合テスト

- 実際のネットワーク通信を使用したエンドツーエンドテスト
- サーバー・クライアント間の通信テスト

### 9.3 テーブル駆動テスト

```go
// テーブル駆動テストの例
func TestNewMessage(t *testing.T) {
    tests := []struct {
        name    string
        content string
        sender  UserID
        wantErr bool
    }{
        {"正常なケース", "こんにちは", "user1", false},
        {"空のメッセージ", "", "user1", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := NewMessage(tt.content, tt.sender)
            if (err != nil) != tt.wantErr {
                t.Errorf("NewMessage() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## 10. 実装計画

実装は以下の順序で進めます：

1. ドメインモデルの実装
2. ネットワークコンポーネントの実装
3. UI コンポーネントの実装
4. アプリケーションサービスの実装
5. メインプログラムの実装
6. テストの実装と実行
7. バグ修正と改善
