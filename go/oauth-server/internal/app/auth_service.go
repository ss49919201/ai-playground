package app

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/ss49919201/ai-playground/go/oauth-server/internal/adapters/storage" // エラー型を参照するため
	"github.com/ss49919201/ai-playground/go/oauth-server/internal/domain"
	"github.com/ss49919201/ai-playground/go/oauth-server/internal/ports"
)

// AuthService は認証と認可に関連するユースケース（主に認可エンドポイント）を処理します。
type AuthService struct {
	clientRepo  ports.ClientRepository
	userRepo    ports.UserRepository
	codeRepo    ports.AuthorizationCodeRepository
	tokenRepo   ports.TokenRepository // インプリシットフローでトークンを発行する場合
	pwHasher    ports.PasswordHasher  // ユーザー認証用
	codeIssuer  ports.CodeIssuer      // 認可コード生成 (副作用)
	tokenIssuer ports.TokenIssuer     // アクセストークン生成 (副作用、インプリシットフロー用)
	clock       ports.Clock           // 時刻取得 (副作用)
	config      AuthServiceConfig     // 認可関連の設定
}

// AuthServiceConfig は AuthService が必要とする設定値を保持します。
type AuthServiceConfig struct {
	AuthCodeLifetime    time.Duration
	AccessTokenLifetime time.Duration // インプリシットフロー用
	// RequirePKCE bool // PKCEを必須とするかどうかのフラグなど
}

// NewAuthService は AuthService の新しいインスタンスを生成します。
func NewAuthService(
	clientRepo ports.ClientRepository,
	userRepo ports.UserRepository,
	codeRepo ports.AuthorizationCodeRepository,
	tokenRepo ports.TokenRepository,
	pwHasher ports.PasswordHasher,
	codeIssuer ports.CodeIssuer,
	tokenIssuer ports.TokenIssuer,
	clock ports.Clock,
	config AuthServiceConfig,
) *AuthService {
	return &AuthService{
		clientRepo:  clientRepo,
		userRepo:    userRepo,
		codeRepo:    codeRepo,
		tokenRepo:   tokenRepo,
		pwHasher:    pwHasher,
		codeIssuer:  codeIssuer,
		tokenIssuer: tokenIssuer,
		clock:       clock,
		config:      config,
	}
}

// AuthorizeRequest は認可エンドポイントへのリクエストパラメータです。
type AuthorizeRequest struct {
	ResponseType string          // "code" または "token"
	ClientID     domain.ClientID // クライアントID
	RedirectURI  string          // リダイレクトURI
	Scope        string          // 要求スコープ (スペース区切り)
	State        string          // CSRF対策のstateパラメータ
	UserID       domain.UserID   // 認証済みユーザーのID (事前に認証が必要)
	// --- PKCE (オプション) ---
	// CodeChallenge       string
	// CodeChallengeMethod string
	// --- ユーザー同意情報 (オプション) ---
	// ConsentGiven bool // ユーザーがスコープに同意したか
	// GrantedScopes []domain.Scope // ユーザーが同意したスコープ
}

// AuthorizeResponse は認可エンドポイントからの成功レスポンスです。
// 実際にはHTTPリダイレクトでパラメータが渡されます。
type AuthorizeResponse struct {
	RedirectURI string // パラメータが付与されたリダイレクト先の完全なURI
	// State       string // 参考情報として含める場合がある
}

// Authorize は認可リクエストを処理し、リダイレクト先のURIまたはエラーを返します。
// ユーザー認証はこのメソッドが呼び出される前に行われている前提です。
// ユーザーのスコープ同意も事前に行われているか、この処理の一部として扱います。
func (s *AuthService) Authorize(ctx context.Context, req AuthorizeRequest) (AuthorizeResponse, error) {
	now := s.clock.Now()

	// 1. クライアント取得と検証
	client, err := s.clientRepo.FindByID(ctx, req.ClientID)
	if err != nil {
		// クライアントが見つからない
		// RFC 6749 Section 4.1.2.1, 4.2.2.1: エラーをリダイレクトURIに返すべきではない
		// エラーページを表示するか、OAuthErrorを返す (HTTP層で処理)
		return AuthorizeResponse{}, NewOAuthError("invalid_client", "指定されたクライアントIDは無効です")
	}

	// 2. リダイレクトURIの検証
	// リクエストで指定された redirect_uri が登録済みのものと一致するか確認
	// redirect_uri がリクエストに含まれていない場合、クライアントに1つだけ登録されていればそれを使用、
	// 複数登録されている場合はエラー (RFC 6749 Section 3.1.2.3)
	var validatedRedirectURI string
	if req.RedirectURI == "" {
		if len(client.RedirectURIs) == 1 {
			validatedRedirectURI = client.RedirectURIs[0]
		} else {
			// エラーページ表示 or OAuthError
			return AuthorizeResponse{}, NewOAuthError("invalid_request", "リダイレクトURIが指定されていないか、複数登録されているクライアントでURIが指定されていません")
		}
	} else {
		if !client.ValidateRedirectURI(req.RedirectURI) {
			// エラーページ表示 or OAuthError
			return AuthorizeResponse{}, NewOAuthError("invalid_request", "登録されていないリダイレクトURIです")
		}
		validatedRedirectURI = req.RedirectURI
	}

	// 3. スコープの検証と決定
	requestedScopes, err := domain.ValidateScope(req.Scope)
	if err != nil {
		// エラーをリダイレクトURIに返す (RFC 6749 Section 4.1.2.1, 4.2.2.1)
		return s.buildErrorRedirect(validatedRedirectURI, "invalid_scope", "無効なスコープ形式です", req.State, req.ResponseType == "token"), nil
	}
	// クライアントが要求スコープを許可されているか
	if !client.ValidateScope(requestedScopes) {
		return s.buildErrorRedirect(validatedRedirectURI, "invalid_scope", "クライアントに許可されていないスコープが含まれています", req.State, req.ResponseType == "token"), nil
	}
	// TODO: ユーザー同意の処理
	// - 過去に同意済みか確認
	// - 同意画面を表示し、ユーザーが許可したスコープを取得 (req.GrantedScopes に設定される想定)
	// - ここでは、リクエストされたスコープが全て同意されたものとする
	grantedScopes := requestedScopes
	if len(grantedScopes) == 0 {
		// 何も同意されなかった場合 (あるいはデフォルトスコープがない場合)
		// return s.buildErrorRedirect(validatedRedirectURI, "access_denied", "ユーザーがアクセスを拒否しました", req.State, req.ResponseType == "token"), nil
	}

	// 4. レスポンスタイプに応じた処理
	responseType := req.ResponseType
	switch responseType {
	case "code": // 認可コードフロー
		// クライアントがこのフローを許可されているか
		if !client.HasGrantType(domain.GrantTypeAuthorizationCode) {
			return s.buildErrorRedirect(validatedRedirectURI, "unauthorized_client", "クライアントは認可コードフローを許可されていません", req.State, false), nil
		}

		// TODO: PKCE チャレンジの検証 (必須の場合)
		// if s.config.RequirePKCE && (req.CodeChallenge == "" || req.CodeChallengeMethod == "") {
		// 	return s.buildErrorRedirect(validatedRedirectURI, "invalid_request", "PKCEコードチャレンジが必要です", req.State, false), nil
		// }

		// 認可コード生成 (副作用)
		codeValue, err := s.codeIssuer.IssueCode()
		if err != nil {
			// TODO: エラーロギング
			return s.buildErrorRedirect(validatedRedirectURI, "server_error", "認可コードの生成に失敗しました", req.State, false), nil
		}
		expiresAt := now.Add(s.config.AuthCodeLifetime)
		authCode, err := domain.NewAuthorizationCode(codeValue, client.ID, req.UserID, validatedRedirectURI, grantedScopes, now, expiresAt /*, req.CodeChallenge, req.CodeChallengeMethod*/)
		if err != nil {
			// TODO: エラーロギング
			return s.buildErrorRedirect(validatedRedirectURI, "server_error", "認可コード情報の生成に失敗しました", req.State, false), nil
		}

		// 認可コード保存 (副作用)
		if err := s.codeRepo.Save(ctx, authCode); err != nil {
			// TODO: エラーロギング
			return s.buildErrorRedirect(validatedRedirectURI, "server_error", "認可コードの保存に失敗しました", req.State, false), nil
		}

		// リダイレクトURIにパラメータを追加
		redirectURL, _ := url.Parse(validatedRedirectURI) // validatedRedirectURI は有効なはず
		query := redirectURL.Query()
		query.Set("code", authCode.Value)
		if req.State != "" {
			query.Set("state", req.State)
		}
		redirectURL.RawQuery = query.Encode()

		return AuthorizeResponse{RedirectURI: redirectURL.String()}, nil

	case "token": // インプリシットフロー (オプション)
		// クライアントがこのフローを許可されているか
		if !client.HasGrantType(domain.GrantTypeImplicit) {
			return s.buildErrorRedirect(validatedRedirectURI, "unauthorized_client", "クライアントはインプリシットフローを許可されていません", req.State, true), nil
		}

		// アクセストークン生成 (副作用)
		accessTokenValue, err := s.tokenIssuer.IssueToken()
		if err != nil {
			// TODO: エラーロギング
			return s.buildErrorRedirect(validatedRedirectURI, "server_error", "アクセストークンの生成に失敗しました", req.State, true), nil
		}
		expiresAt := now.Add(s.config.AccessTokenLifetime)
		accessToken, err := domain.NewToken(accessTokenValue, domain.TokenTypeBearer, client.ID, req.UserID, grantedScopes, now, expiresAt)
		if err != nil {
			// TODO: エラーロギング
			return s.buildErrorRedirect(validatedRedirectURI, "server_error", "アクセストークン情報の生成に失敗しました", req.State, true), nil
		}

		// トークン保存 (オプション) - インプリシットフローでは通常保存しない
		// if err := s.tokenRepo.Save(ctx, accessToken); err != nil { ... }

		// リダイレクトURIのフラグメントにパラメータを追加 (RFC 6749 Section 4.2.2)
		redirectURL, _ := url.Parse(validatedRedirectURI)
		fragment := url.Values{}
		fragment.Set("access_token", accessToken.Value)
		fragment.Set("token_type", string(accessToken.Type))
		fragment.Set("expires_in", fmt.Sprintf("%d", int(s.config.AccessTokenLifetime.Seconds())))
		fragment.Set("scope", domain.FormatScopes(grantedScopes))
		if req.State != "" {
			fragment.Set("state", req.State)
		}
		redirectURL.Fragment = fragment.Encode() // クエリではなくフラグメント

		return AuthorizeResponse{RedirectURI: redirectURL.String()}, nil

	default:
		// サポートされていないレスポンスタイプ
		return s.buildErrorRedirect(validatedRedirectURI, "unsupported_response_type", fmt.Sprintf("サポートされていないレスポンスタイプです: %s", responseType), req.State, false), nil // isImplicit はどちらでも良い
	}
}

// AuthenticateUser はユーザー名とパスワードでユーザーを認証します。
// これは認可エンドポイント前のログイン処理などで使用されることを想定しています。
// TokenService の authenticateUser と重複しますが、責務が若干異なる可能性があります。
func (s *AuthService) AuthenticateUser(ctx context.Context, username, password string) (domain.User, error) {
	if username == "" || password == "" {
		return domain.User{}, errors.New("ユーザー名とパスワードは必須です") // OAuthErrorではない通常のエラー
	}
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return domain.User{}, errors.New("ユーザー名またはパスワードが無効です")
		}
		// TODO: エラーロギング
		return domain.User{}, errors.New("ユーザー認証中にエラーが発生しました")
	}
	match, err := s.pwHasher.Compare(user.HashedPassword, password)
	if err != nil {
		// ハッシュ比較エラー
		// TODO: エラーロギング
		return domain.User{}, errors.New("ユーザー認証中にエラーが発生しました")
	}
	if !match {
		// パスワード不一致
		return domain.User{}, errors.New("ユーザー名またはパスワードが無効です")
	}
	return user, nil
}

// buildErrorRedirect は認可エンドポイントでのエラー時にリダイレクトするURIを構築します。
func (s *AuthService) buildErrorRedirect(redirectURI, errorCode, errorDesc, state string, isImplicit bool) AuthorizeResponse {
	redirectURL, err := url.Parse(redirectURI)
	if err != nil {
		// リダイレクトURIが無効な場合はリダイレクトしない方が安全
		// 代わりにエラーページを表示するなどの処理が必要
		// ここでは仮に空のレスポンスを返す (呼び出し元でエラー処理が必要)
		return AuthorizeResponse{}
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
	return AuthorizeResponse{RedirectURI: redirectURL.String()}
}
