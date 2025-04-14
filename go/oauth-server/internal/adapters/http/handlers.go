package httpadapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ss49919201/ai-kata/go/oauth-server/internal/adapters/storage"
	"github.com/ss49919201/ai-kata/go/oauth-server/internal/app"
	"github.com/ss49919201/ai-kata/go/oauth-server/internal/domain"
)

// handleAuthorize は認可エンドポイント (`/oauth/authorize`) のリクエストを処理します。
// GET または POST リクエストを受け付けます。
// ユーザー認証と同意が事前に完了していることを前提とします。
func (s *Server) handleAuthorize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		s.renderErrorPage(w, r, http.StatusMethodNotAllowed, "Method Not Allowed", "GET または POST メソッドを使用してください。")
		return
	}

	// TODO: ユーザー認証状態の確認
	// セッションやクッキーから認証済みユーザーの情報を取得する。
	// 認証されていない場合はログインページにリダイレクトする。
	// この例では、認証済みで UserID が取得できていると仮定します。
	userID := domain.UserID("dummy-user-id-123") // 仮のユーザーID

	// リクエストパラメータの取得
	query := r.URL.Query()
	req := app.AuthorizeRequest{
		ResponseType: query.Get("response_type"),
		ClientID:     domain.ClientID(query.Get("client_id")),
		RedirectURI:  query.Get("redirect_uri"), // 省略される可能性あり
		Scope:        query.Get("scope"),
		State:        query.Get("state"),
		UserID:       userID,
		// TODO: PKCE パラメータ (code_challenge, code_challenge_method) の取得
	}

	// 必須パラメータのチェック
	if req.ResponseType == "" || req.ClientID == "" {
		// RFC 6749 Section 4.1.2.1: redirect_uri が無効な場合を除き、エラーをリダイレクトしない
		s.renderErrorPage(w, r, http.StatusBadRequest, "invalid_request", "response_type と client_id は必須パラメータです。")
		return
	}

	// アプリケーションサービスの呼び出し
	resp, err := s.authService.Authorize(r.Context(), req)

	// エラーハンドリング
	if err != nil {
		var oauthErr *app.OAuthError
		if errors.As(err, &oauthErr) {
			// AuthService が OAuthError を返した場合 (リダイレクトすべきでないエラー)
			s.renderErrorPage(w, r, http.StatusBadRequest, oauthErr.Code, oauthErr.Description)
		} else {
			// その他の予期せぬエラー
			// TODO: エラーロギング
			s.renderErrorPage(w, r, http.StatusInternalServerError, "server_error", "認可処理中に内部エラーが発生しました。")
		}
		return
	}

	// Authorize がエラーを返さなかった場合、レスポンスにはリダイレクトURIが含まれているはず
	// (エラーの場合でもリダイレクトするケースは Authorize 内で処理されている)
	if resp.RedirectURI == "" {
		// Authorize がリダイレクトURIを返さなかった場合 (通常はエラーケース)
		// TODO: エラーロギング
		s.renderErrorPage(w, r, http.StatusInternalServerError, "server_error", "リダイレクト先の生成に失敗しました。")
		return
	}

	// 成功レスポンス (HTTPリダイレクト)
	// RFC 6749 Section 3.1.2.5: キャッシュを防ぐヘッダーを追加することが推奨される
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	http.Redirect(w, r, resp.RedirectURI, http.StatusFound) // 302 Found
}

// handleToken はトークンエンドポイント (`/oauth/token`) のリクエストを処理します。
// POST リクエストのみを受け付けます。
func (s *Server) handleToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.renderJSONError(w, http.StatusMethodNotAllowed, "invalid_request", "POST メソッドを使用してください。")
		return
	}

	// Content-Type のチェック (RFC 6749 Section 2.3.1)
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		s.renderJSONError(w, http.StatusUnsupportedMediaType, "invalid_request", "Content-Type は application/x-www-form-urlencoded である必要があります。")
		return
	}

	// クライアント認証 (RFC 6749 Section 2.3.1)
	// Basic 認証ヘッダー または リクエストボディの client_id/client_secret
	clientIDStr, clientSecret, ok := r.BasicAuth()
	if !ok {
		// Basic 認証がない場合はリクエストボディから取得試行
		if err := r.ParseForm(); err != nil {
			s.renderJSONError(w, http.StatusBadRequest, "invalid_request", "リクエストボディの解析に失敗しました。")
			return
		}
		clientIDStr = r.PostFormValue("client_id")
		clientSecret = r.PostFormValue("client_secret")

		// client_id は必須
		if clientIDStr == "" {
			s.renderJSONError(w, http.StatusBadRequest, "invalid_client", "クライアント認証情報 (client_id) が必要です。")
			return
		}
		// client_secret は Confidential Client の場合に必要だが、
		// Public Client かどうかはこの時点では不明なため、TokenService内で検証する
	}

	// リクエストボディのパース (Basic認証の場合も必要)
	if err := r.ParseForm(); err != nil && r.ContentLength > 0 { // ボディがある場合のみパース
		s.renderJSONError(w, http.StatusBadRequest, "invalid_request", "リクエストボディの解析に失敗しました。")
		return
	}

	// アプリケーションサービスへのリクエストを作成
	req := app.IssueTokenRequest{
		GrantType:    r.PostFormValue("grant_type"),
		Code:         r.PostFormValue("code"),
		RedirectURI:  r.PostFormValue("redirect_uri"),
		ClientID:     domain.ClientID(clientIDStr),
		ClientSecret: clientSecret,
		Username:     r.PostFormValue("username"),
		Password:     r.PostFormValue("password"),
		RefreshToken: r.PostFormValue("refresh_token"),
		Scope:        r.PostFormValue("scope"),
		// TODO: PKCE パラメータ (code_verifier) の取得
	}

	// GrantType は必須
	if req.GrantType == "" {
		s.renderJSONError(w, http.StatusBadRequest, "invalid_request", "grant_type パラメータは必須です。")
		return
	}

	// アプリケーションサービスの呼び出し
	resp, err := s.tokenService.IssueToken(r.Context(), req)

	// エラーハンドリング
	if err != nil {
		var oauthErr *app.OAuthError
		if errors.As(err, &oauthErr) {
			// TokenService が OAuthError を返した場合
			statusCode := http.StatusBadRequest // デフォルト
			if oauthErr.Code == "invalid_client" {
				statusCode = http.StatusUnauthorized
			} else if oauthErr.Code == "server_error" {
				statusCode = http.StatusInternalServerError
			}
			s.renderJSONError(w, statusCode, oauthErr.Code, oauthErr.Description)
		} else {
			// その他の予期せぬエラー
			// TODO: エラーロギング
			s.renderJSONError(w, http.StatusInternalServerError, "server_error", "トークン発行処理中に内部エラーが発生しました。")
		}
		return
	}

	// 成功レスポンス (JSON)
	// RFC 6749 Section 5.1: キャッシュを防ぐヘッダーを追加
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		// JSONエンコードエラー (通常は発生しないはず)
		// TODO: エラーロギング
	}
}

// handleIntrospect はトークンイントロスペクションエンドポイント (`/oauth/introspect`) を処理します。
// RFC 7662 準拠。POST リクエストのみを受け付けます。
// リソースサーバーからのリクエストを想定し、認証が必要です (例: Basic認証、Bearerトークン)。
func (s *Server) handleIntrospect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.renderJSONError(w, http.StatusMethodNotAllowed, "invalid_request", "POST メソッドを使用してください。")
		return
	}

	// TODO: リソースサーバーの認証
	// - Basic認証 (クライアントID/シークレット)
	// - Bearerトークン (特別な権限を持つトークン)
	// 認証されていない場合は 401 Unauthorized を返す

	// リクエストボディのパース (application/x-www-form-urlencoded)
	if err := r.ParseForm(); err != nil {
		s.renderJSONError(w, http.StatusBadRequest, "invalid_request", "リクエストボディの解析に失敗しました。")
		return
	}
	tokenValue := r.PostFormValue("token")
	// tokenTypeHint := r.PostFormValue("token_type_hint") // オプション

	if tokenValue == "" {
		s.renderJSONError(w, http.StatusBadRequest, "invalid_request", "token パラメータは必須です。")
		return
	}

	// アプリケーションサービスの呼び出し
	resp, err := s.tokenService.ValidateToken(r.Context(), tokenValue)
	if err != nil {
		// ValidateToken は通常エラーを返さない設計だが、予期せぬエラーの場合
		// TODO: エラーロギング
		s.renderJSONError(w, http.StatusInternalServerError, "server_error", "トークン検証中に内部エラーが発生しました。")
		return
	}

	// 成功レスポンス (JSON)
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		// TODO: エラーロギング
	}
}

// handleRevoke はトークン失効エンドポイント (`/oauth/revoke`) を処理します。
// RFC 7009 準拠。POST リクエストのみを受け付けます。
// クライアント認証が必要です。
func (s *Server) handleRevoke(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.renderJSONError(w, http.StatusMethodNotAllowed, "invalid_request", "POST メソッドを使用してください。")
		return
	}

	// クライアント認証 (Basic または ボディ)
	clientIDStr, clientSecret, ok := r.BasicAuth()
	if !ok {
		if err := r.ParseForm(); err != nil {
			s.renderJSONError(w, http.StatusBadRequest, "invalid_request", "リクエストボディの解析に失敗しました。")
			return
		}
		clientIDStr = r.PostFormValue("client_id")
		clientSecret = r.PostFormValue("client_secret")
		if clientIDStr == "" { // シークレットは Public Client の場合省略可能
			s.renderJSONError(w, http.StatusBadRequest, "invalid_client", "クライアント認証情報 (client_id) が必要です。")
			return
		}
	} else {
		// Basic認証の場合もボディをパースする必要がある
		if err := r.ParseForm(); err != nil && r.ContentLength > 0 {
			s.renderJSONError(w, http.StatusBadRequest, "invalid_request", "リクエストボディの解析に失敗しました。")
			return
		}
	}
	clientID := domain.ClientID(clientIDStr)

	// 失効させるトークンを取得
	tokenValue := r.PostFormValue("token")
	// tokenTypeHint := r.PostFormValue("token_type_hint") // オプション

	if tokenValue == "" {
		s.renderJSONError(w, http.StatusBadRequest, "invalid_request", "token パラメータは必須です。")
		return
	}

	// まずクライアントを認証
	_, err := s.clientService.AuthenticateClient(r.Context(), clientID, clientSecret)
	if err != nil {
		var oauthErr *app.OAuthError
		// clientService.AuthenticateClient は OAuthError を返さない設計なので、エラーの種類を判別
		// if errors.Is(err, storage.ErrClientNotFound) || strings.Contains(err.Error(), "invalid_client") { // 仮のエラーチェック
		if errors.Is(err, storage.ErrClientNotFound) || (err != nil && err.Error() == "クライアント認証に失敗しました (invalid_client)") { // client_serviceのエラーメッセージに合わせる
			statusCode := http.StatusBadRequest
			if oauthErr.Code == "invalid_client" {
				statusCode = http.StatusUnauthorized
			}
			s.renderJSONError(w, statusCode, oauthErr.Code, oauthErr.Description)
		} else {
			// TODO: エラーロギング
			s.renderJSONError(w, http.StatusInternalServerError, "server_error", "クライアント認証中にエラーが発生しました。")
		}
		return
	}

	// アプリケーションサービスの呼び出し
	err = s.tokenService.RevokeToken(r.Context(), tokenValue, clientID)
	if err != nil {
		var oauthErr *app.OAuthError
		if errors.As(err, &oauthErr) {
			// 通常、RevokeToken は RFC 7009 に従い、トークンが見つからなくてもエラーを返さない想定
			// ここでエラーが返るのは server_error など
			s.renderJSONError(w, http.StatusInternalServerError, oauthErr.Code, oauthErr.Description)
		} else {
			// TODO: エラーロギング
			s.renderJSONError(w, http.StatusInternalServerError, "server_error", "トークン失効処理中に内部エラーが発生しました。")
		}
		return
	}

	// 成功レスポンス (RFC 7009 Section 2.2: 200 OK、ボディは空)
	w.WriteHeader(http.StatusOK)
}

// handleClients はクライアント管理エンドポイント (`/oauth/clients`) を処理します。
// これは管理用のAPIであり、適切な認証/認可が必要です。
func (s *Server) handleClients(w http.ResponseWriter, r *http.Request) {
	// TODO: 管理者認証/認可の実装
	// 例: 特定のIPアドレスからのみ許可、管理用トークン、Basic認証など
	// if !isAdmin(r) {
	// 	http.Error(w, "Forbidden", http.StatusForbidden)
	// 	return
	// }

	switch r.Method {
	case http.MethodPost: // クライアント登録
		s.handleRegisterClient(w, r)
	case http.MethodGet: // クライアント取得 (ID指定 or 一覧)
		s.handleGetClient(w, r)
	// case http.MethodPut: // クライアント更新
	// 	s.handleUpdateClient(w, r)
	// case http.MethodDelete: // クライアント削除
	// 	s.handleDeleteClient(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// handleRegisterClient はクライアント登録リクエストを処理します。
func (s *Server) handleRegisterClient(w http.ResponseWriter, r *http.Request) {
	var req app.RegisterClientRequest
	// Content-Type が application/json であることを期待
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.renderJSONError(w, http.StatusBadRequest, "invalid_request", "リクエストボディ(JSON)の解析に失敗しました。")
		return
	}

	// アプリケーションサービスの呼び出し
	resp, err := s.clientService.RegisterClient(r.Context(), req)
	if err != nil {
		// ドメインバリデーションエラーなども含まれる可能性
		// TODO: エラーの種類に応じてステータスコードを分ける (例: Bad Request)
		// TODO: エラーロギング
		s.renderJSONError(w, http.StatusInternalServerError, "server_error", fmt.Sprintf("クライアント登録に失敗しました: %v", err))
		return
	}

	// 成功レスポンス (JSON)
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusCreated) // 201 Created
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		// TODO: エラーロギング
	}
}

// handleGetClient はクライアント情報取得リクエストを処理します。
func (s *Server) handleGetClient(w http.ResponseWriter, r *http.Request) {
	// パスから ClientID を取得 (例: /oauth/clients/{client_id})
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	var clientID domain.ClientID
	if len(pathParts) == 3 { // ["oauth", "clients", "{client_id}"]
		clientID = domain.ClientID(pathParts[2])
	} else if len(pathParts) == 2 { // ["oauth", "clients"] -> 一覧取得 (未実装)
		http.Error(w, "Not Implemented: Client listing", http.StatusNotImplemented)
		return
	} else {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	// アプリケーションサービスの呼び出し
	resp, err := s.clientService.GetClient(r.Context(), clientID)
	if err != nil {
		if errors.Is(err, storage.ErrClientNotFound) {
			http.Error(w, "Client Not Found", http.StatusNotFound)
		} else {
			// TODO: エラーロギング
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// 成功レスポンス (JSON)
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		// TODO: エラーロギング
	}
}

// --- ヘルパー関数 ---

// renderJSONError は OAuth 2.0 形式のエラーレスポンスをJSONで返します。
func (s *Server) renderJSONError(w http.ResponseWriter, statusCode int, errorCode, errorDesc string) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	// トークンエンドポイントのエラーレスポンスではキャッシュ制御ヘッダーが必要 (RFC 6749 Section 5.2)
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.WriteHeader(statusCode)
	errResp := app.OAuthError{Code: errorCode, Description: errorDesc}
	if err := json.NewEncoder(w).Encode(errResp); err != nil {
		// JSONエンコードエラー時のフォールバック
		// TODO: エラーロギング
		http.Error(w, `{"error":"server_error","error_description":"Failed to encode error response"}`, http.StatusInternalServerError)
	}
}

// renderErrorPage はユーザー向けのエラーページを表示します (HTMLなど)。
// 認可エンドポイントでリダイレクトできない場合などに使用します。
func (s *Server) renderErrorPage(w http.ResponseWriter, r *http.Request, statusCode int, errorCode, errorDesc string) {
	// TODO: 実際のエラーページテンプレートをレンダリングする
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, "<html><body><h1>エラー (%s)</h1><p>%s</p></body></html>", errorCode, errorDesc)
}

// buildErrorRedirect は認可エンドポイントでのエラー時にリダイレクトするURIを構築します。
// AuthService にも同様のメソッドがあるため、共通化を検討。
func (s *Server) buildErrorRedirect(redirectURI, errorCode, errorDesc, state string, isImplicit bool) (string, error) {
	redirectURL, err := url.Parse(redirectURI)
	if err != nil {
		return "", fmt.Errorf("無効なリダイレクトURIです: %w", err)
	}

	if isImplicit { // インプリシットフローはフラグメント
		fragment := url.Values{}
		fragment.Set("error", errorCode)
		if errorDesc != "" {
			fragment.Set("error_description", errorDesc)
		}
		if state != "" {
			fragment.Set("state", state)
		}
		redirectURL.Fragment = fragment.Encode()
	} else { // 認可コードフローはクエリパラメータ
		query := redirectURL.Query()
		query.Set("error", errorCode)
		if errorDesc != "" {
			query.Set("error_description", errorDesc)
		}
		if state != "" {
			query.Set("state", state)
		}
		redirectURL.RawQuery = query.Encode()
	}
	return redirectURL.String(), nil
}
