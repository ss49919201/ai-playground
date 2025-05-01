package httpadapter

import (
	"net/http"

	"github.com/ss49919201/ai-kata/go/oauth-server/internal/app"
)

// Server はHTTPサーバーの依存関係とルーターを保持します。
// http.Handler インターフェースを実装します。
type Server struct {
	authService   *app.AuthService
	tokenService  *app.TokenService
	clientService *app.ClientService
	mux           *http.ServeMux // または他のルーター (chi, gorilla/mux など)
	// logger        *log.Logger    // ロガーなど、他の依存関係も追加可能
}

// NewServer はHTTPサーバーの新しいインスタンスを生成し、
// 依存関係を注入してハンドラーを登録します。
func NewServer(
	authSvc *app.AuthService,
	tokenSvc *app.TokenService,
	clientSvc *app.ClientService,
) *Server {
	s := &Server{
		authService:   authSvc,
		tokenService:  tokenSvc,
		clientService: clientSvc,
		mux:           http.NewServeMux(), // 標準のServeMuxを使用
	}
	s.registerHandlers() // ハンドラーをmuxに登録
	return s
}

// ServeHTTP は http.Handler インターフェースを実装し、
// 受け取ったリクエストを内部のルーター (mux) に委譲します。
// ここでリクエスト全体に適用するミドルウェア（ロギング、リカバリなど）を実装することも可能です。
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 例: リクエストロギングミドルウェア
	// startTime := time.Now()
	// s.logger.Printf("Started %s %s", r.Method, r.URL.Path)
	// defer func() {
	// 	s.logger.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(startTime))
	// }()

	// ルーターに処理を委譲
	s.mux.ServeHTTP(w, r)
}

// registerHandlers はサーバーの各エンドポイントに対応するハンドラー関数を
// ルーター (mux) に登録します。
func (s *Server) registerHandlers() {
	// OAuth 2.0 エンドポイント
	s.mux.HandleFunc("/oauth/authorize", s.handleAuthorize) // 認可エンドポイント
	s.mux.HandleFunc("/oauth/token", s.handleToken)         // トークンエンドポイント

	// オプションのエンドポイント (RFC 7662, RFC 7009)
	s.mux.HandleFunc("/oauth/introspect", s.handleIntrospect) // トークンイントロスペクション
	s.mux.HandleFunc("/oauth/revoke", s.handleRevoke)         // トークン失効

	// 管理用エンドポイント
	// 注意: これらのエンドポイントは適切な認証/認可で保護する必要があります。
	s.mux.HandleFunc("/oauth/clients", s.handleClients)  // クライアント一覧取得・登録
	s.mux.HandleFunc("/oauth/clients/", s.handleClients) // 特定クライアント取得・更新・削除 (パスでIDを指定)

	// TODO: ユーザー認証ページ、同意ページなどのハンドラーも必要に応じて追加
	// s.mux.HandleFunc("/login", s.handleLogin)
	// s.mux.HandleFunc("/consent", s.handleConsent)
}
