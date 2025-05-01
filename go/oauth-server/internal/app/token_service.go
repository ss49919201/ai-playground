package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ss49919201/ai-playground/go/oauth-server/internal/adapters/storage" // エラー型を参照するため
	"github.com/ss49919201/ai-playground/go/oauth-server/internal/domain"
	"github.com/ss49919201/ai-playground/go/oauth-server/internal/ports"
)

// TokenService はトークンの発行、検証、失効に関連するユースケースを処理します。
type TokenService struct {
	clientRepo  ports.ClientRepository
	userRepo    ports.UserRepository
	codeRepo    ports.AuthorizationCodeRepository
	tokenRepo   ports.TokenRepository
	pwHasher    ports.PasswordHasher // ユーザー認証 (Password Grant) やクライアント認証で使用
	tokenIssuer ports.TokenIssuer    // トークン生成 (副作用)
	clock       ports.Clock          // 時刻取得 (副作用)
	config      TokenServiceConfig   // トークン関連の設定
}

// TokenServiceConfig は TokenService が必要とする設定値を保持します。
// config.Config から必要な値を取り出して渡します。
type TokenServiceConfig struct {
	AccessTokenLifetime  time.Duration
	RefreshTokenLifetime time.Duration
	// IssueRefreshToken bool // リフレッシュトークンを発行するかどうかのフラグなど
}

// NewTokenService は TokenService の新しいインスタンスを生成します。
func NewTokenService(
	clientRepo ports.ClientRepository,
	userRepo ports.UserRepository,
	codeRepo ports.AuthorizationCodeRepository,
	tokenRepo ports.TokenRepository,
	pwHasher ports.PasswordHasher,
	tokenIssuer ports.TokenIssuer,
	clock ports.Clock,
	config TokenServiceConfig,
) *TokenService {
	return &TokenService{
		clientRepo:  clientRepo,
		userRepo:    userRepo,
		codeRepo:    codeRepo,
		tokenRepo:   tokenRepo,
		pwHasher:    pwHasher,
		tokenIssuer: tokenIssuer,
		clock:       clock,
		config:      config,
	}
}

// IssueTokenRequest はトークン発行リクエストのパラメータです。
// 各 Grant Type で使用されるフィールドが異なります。
type IssueTokenRequest struct {
	GrantType    string // "authorization_code", "password", "client_credentials", "refresh_token"
	Code         string // GrantType: "authorization_code"
	RedirectURI  string // GrantType: "authorization_code"
	ClientID     domain.ClientID
	ClientSecret string // Confidential Client の認証用
	Username     string // GrantType: "password"
	Password     string // GrantType: "password"
	RefreshToken string // GrantType: "refresh_token"
	Scope        string // オプション: 要求するスコープ (スペース区切り)
}

// IssueTokenResponse はトークン発行成功時のレスポンスパラメータです。
// RFC 6749 Section 5.1 準拠。
type IssueTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`              // 通常は "Bearer"
	ExpiresIn    int    `json:"expires_in"`              // アクセストークンの有効期間 (秒)
	RefreshToken string `json:"refresh_token,omitempty"` // 発行された場合のみ
	Scope        string `json:"scope,omitempty"`         // 実際に許可されたスコープ (スペース区切り)
}

// OAuthError はトークンエンドポイントでのエラーレスポンスです。
// RFC 6749 Section 5.2 準拠。
type OAuthError struct {
	Code        string `json:"error"`                       // エラーコード (例: "invalid_request")
	Description string `json:"error_description,omitempty"` // エラーの詳細説明 (オプション)
	URI         string `json:"error_uri,omitempty"`         // エラーに関する詳細情報URI (オプション)
	// HTTPStatusCode int    `json:"-"`                       // 対応するHTTPステータスコード
}

// Error は error インターフェースを実装します。
func (e *OAuthError) Error() string {
	return fmt.Sprintf("OAuth Error: %s (%s)", e.Code, e.Description)
}

// NewOAuthError は OAuthError を生成します。
func NewOAuthError(code, description string) *OAuthError {
	return &OAuthError{Code: code, Description: description}
}

// IssueToken はトークン発行リクエストを処理し、トークンまたはエラーを返します。
func (s *TokenService) IssueToken(ctx context.Context, req IssueTokenRequest) (IssueTokenResponse, error) {
	now := s.clock.Now()

	// 1. クライアント認証
	// Public Client の場合は ClientSecret は空文字列になる
	client, err := s.authenticateClient(ctx, req.ClientID, req.ClientSecret)
	if err != nil {
		// authenticateClient が返すエラーは既に OAuthError 形式のはず
		return IssueTokenResponse{}, err
	}

	// 2. Grant Type の検証と処理
	var userID domain.UserID
	var grantedScopes []domain.Scope
	var originalRefreshTokenScopes []domain.Scope // リフレッシュトークンフロー用

	grantType := domain.GrantType(req.GrantType)

	// クライアントがこの Grant Type を許可されているかチェック
	if !client.HasGrantType(grantType) {
		// refresh_token は他のフローの結果として許可されるため、個別のチェックは不要かもしれない
		if grantType != domain.GrantTypeRefreshToken {
			return IssueTokenResponse{}, NewOAuthError("unauthorized_client", fmt.Sprintf("クライアントはこの認可フロー (%s) を許可されていません", grantType))
		}
		// リフレッシュトークン自体が許可されているかは client.HasGrantType(domain.GrantTypeRefreshToken) でチェックする
		if !client.HasGrantType(domain.GrantTypeRefreshToken) {
			return IssueTokenResponse{}, NewOAuthError("unauthorized_client", "クライアントはリフレッシュトークンフローを許可されていません")
		}
	}

	switch grantType {
	case domain.GrantTypeAuthorizationCode:
		// 認可コードの検証
		authCode, err := s.validateAuthorizationCode(ctx, req.Code, client.ID, req.RedirectURI, now)
		if err != nil {
			return IssueTokenResponse{}, err // validateAuthorizationCode が OAuthError を返す
		}
		userID = authCode.UserID
		grantedScopes = authCode.Scopes
		// 認可コードを削除 (一度きり有効)
		if err := s.codeRepo.Delete(ctx, authCode.Value); err != nil {
			// 削除失敗はログに残すが、トークン発行は続行する (べきか？)
			// TODO: エラーロギング
		}

	case domain.GrantTypePassword:
		// ユーザー認証
		user, err := s.authenticateUser(ctx, req.Username, req.Password)
		if err != nil {
			return IssueTokenResponse{}, err // authenticateUser が OAuthError を返す
		}
		userID = user.ID
		// スコープの決定
		requestedScopes, err := domain.ValidateScope(req.Scope)
		if err != nil {
			return IssueTokenResponse{}, NewOAuthError("invalid_scope", "無効なスコープ形式です")
		}
		grantedScopes = s.determineGrantedScopes(client.Scopes, nil, requestedScopes) // TODO: ユーザー固有スコープも考慮

	case domain.GrantTypeClientCredentials:
		// クライアント自身のトークンなので UserID はなし
		userID = ""
		// スコープの決定
		requestedScopes, err := domain.ValidateScope(req.Scope)
		if err != nil {
			return IssueTokenResponse{}, NewOAuthError("invalid_scope", "無効なスコープ形式です")
		}
		grantedScopes = s.determineGrantedScopes(client.Scopes, nil, requestedScopes)

	case domain.GrantTypeRefreshToken:
		// リフレッシュトークンの検証
		refreshToken, err := s.validateRefreshToken(ctx, req.RefreshToken, client.ID, now)
		if err != nil {
			return IssueTokenResponse{}, err // validateRefreshToken が OAuthError を返す
		}
		userID = refreshToken.UserID
		originalRefreshTokenScopes = refreshToken.Scopes // 元のスコープを保持

		// オプション: リフレッシュトークンローテーション - 古いトークンを失効させる
		// if err := s.tokenRepo.Delete(ctx, refreshToken.Value); err != nil {
		// 	// ログに残すが処理は続行
		// 	// TODO: エラーロギング
		// }

		// スコープの決定 (リクエストされたスコープが元のスコープのサブセットであることを確認)
		requestedScopes, err := domain.ValidateScope(req.Scope)
		if err != nil {
			return IssueTokenResponse{}, NewOAuthError("invalid_scope", "無効なスコープ形式です")
		}
		if len(requestedScopes) == 0 {
			// リクエストでスコープが指定されなかった場合は、元のスコープを引き継ぐ
			grantedScopes = originalRefreshTokenScopes
		} else {
			// リクエストされたスコープが元のスコープに含まれているかチェック
			tempToken := domain.Token{Scopes: originalRefreshTokenScopes} // チェック用の一時トークン
			if !tempToken.HasAllScopes(requestedScopes) {
				return IssueTokenResponse{}, NewOAuthError("invalid_scope", "要求されたスコープが元のリフレッシュトークンのスコープを超えています")
			}
			grantedScopes = requestedScopes // 要求されたスコープを許可
		}

	default:
		return IssueTokenResponse{}, NewOAuthError("unsupported_grant_type", fmt.Sprintf("サポートされていないGrant Typeです: %s", grantType))
	}

	// 3. アクセストークン生成
	accessTokenValue, err := s.tokenIssuer.IssueToken()
	if err != nil {
		// TODO: エラーロギング
		return IssueTokenResponse{}, NewOAuthError("server_error", "アクセストークンの生成に失敗しました")
	}
	accessTokenExpiresAt := now.Add(s.config.AccessTokenLifetime)
	accessToken, err := domain.NewToken(accessTokenValue, domain.TokenTypeBearer, client.ID, userID, grantedScopes, now, accessTokenExpiresAt)
	if err != nil {
		// ドメインレベルのエラー
		// TODO: エラーロギング
		return IssueTokenResponse{}, NewOAuthError("server_error", "アクセストークン情報の生成に失敗しました")
	}
	if err := s.tokenRepo.Save(ctx, accessToken); err != nil {
		// TODO: エラーロギング
		return IssueTokenResponse{}, NewOAuthError("server_error", "アクセストークンの保存に失敗しました")
	}

	// 4. リフレッシュトークン生成 (必要な場合)
	var refreshTokenValue string
	// リフレッシュトークンを発行する条件:
	// - クライアントが refresh_token grant type を許可されている
	// - 今回のフローが Authorization Code または Password または Refresh Token である
	issueRefreshToken := client.HasGrantType(domain.GrantTypeRefreshToken) &&
		(grantType == domain.GrantTypeAuthorizationCode || grantType == domain.GrantTypePassword || grantType == domain.GrantTypeRefreshToken)

	if issueRefreshToken {
		// オプション: リフレッシュトークンローテーションの場合、新しい値を生成
		// そうでない場合、既存のリフレッシュトークンを再利用するかどうかはポリシーによる
		// ここでは常に新しい値を生成する
		newRefreshTokenValue, err := s.tokenIssuer.IssueToken()
		if err != nil {
			// TODO: エラーロギング
			return IssueTokenResponse{}, NewOAuthError("server_error", "リフレッシュトークンの生成に失敗しました")
		}
		refreshTokenValue = newRefreshTokenValue

		refreshTokenExpiresAt := now.Add(s.config.RefreshTokenLifetime)
		// リフレッシュトークンに含めるスコープは、今回許可されたスコープと同じにする
		refreshToken, err := domain.NewToken(refreshTokenValue, "", client.ID, userID, grantedScopes, now, refreshTokenExpiresAt) // Typeは空など
		if err != nil {
			// TODO: エラーロギング
			return IssueTokenResponse{}, NewOAuthError("server_error", "リフレッシュトークン情報の生成に失敗しました")
		}
		if err := s.tokenRepo.Save(ctx, refreshToken); err != nil {
			// TODO: エラーロギング
			return IssueTokenResponse{}, NewOAuthError("server_error", "リフレッシュトークンの保存に失敗しました")
		}
	}

	// 5. レスポンス生成
	resp := IssueTokenResponse{
		AccessToken:  accessToken.Value,
		TokenType:    string(accessToken.Type),
		ExpiresIn:    int(s.config.AccessTokenLifetime.Seconds()),
		RefreshToken: refreshTokenValue, // 生成した場合のみ設定
		Scope:        domain.FormatScopes(grantedScopes),
	}

	return resp, nil
}

// ValidateTokenResponse はトークン検証 (Introspection) のレスポンスです。
// RFC 7662 Section 2.2 準拠。
type ValidateTokenResponse struct {
	Active    bool     `json:"active"`               // トークンが有効かどうか
	Scope     string   `json:"scope,omitempty"`      // 許可されたスコープ (スペース区切り)
	ClientID  string   `json:"client_id,omitempty"`  // クライアントID
	Username  string   `json:"username,omitempty"`   // ユーザー名 (ユーザーに紐づく場合)
	TokenType string   `json:"token_type,omitempty"` // トークンタイプ (例: "Bearer")
	ExpiresAt int64    `json:"exp,omitempty"`        // 有効期限 (Unixタイムスタンプ)
	IssuedAt  int64    `json:"iat,omitempty"`        // 発行日時 (Unixタイムスタンプ)
	Subject   string   `json:"sub,omitempty"`        // 主体 (ユーザーIDなど)
	Audience  []string `json:"aud,omitempty"`        // 対象者 (クライアントIDなど)
	Issuer    string   `json:"iss,omitempty"`        // 発行者
	JwtID     string   `json:"jti,omitempty"`        // JWT ID
}

// ValidateToken は提供されたトークン文字列を検証します。
// アクセストークンまたはリフレッシュトークンの可能性があります。
// 有効な場合はトークン情報を含むレスポンスを、無効な場合は active: false のレスポンスを返します。
func (s *TokenService) ValidateToken(ctx context.Context, tokenValue string) (ValidateTokenResponse, error) {
	token, err := s.tokenRepo.FindByValue(ctx, tokenValue)
	now := s.clock.Now()

	// トークンが見つからない、または有効期限切れの場合
	if err != nil || token.IsExpired(now) {
		// RFC 7662 では、無効なトークンの場合でもエラーではなく active: false を返す
		return ValidateTokenResponse{Active: false}, nil
	}

	// 有効なトークン情報をレスポンスに設定
	resp := ValidateTokenResponse{
		Active:    true,
		Scope:     domain.FormatScopes(token.Scopes),
		ClientID:  string(token.ClientID),
		TokenType: string(token.Type), // Bearer または空 (Refresh Token)
		ExpiresAt: token.ExpiresAt.Unix(),
		IssuedAt:  token.IssuedAt.Unix(),
		Subject:   string(token.UserID),             // UserID を Subject とする
		Audience:  []string{string(token.ClientID)}, // ClientID を Audience とする
		// Issuer:    s.config.Issuer, // 設定から取得
		// JwtID:     token.JwtID, // JWTの場合
	}

	// ユーザー名を取得 (オプション)
	if token.UserID != "" {
		user, err := s.userRepo.FindByID(ctx, token.UserID)
		if err == nil {
			resp.Username = user.Username
		} else {
			// ユーザーが見つからない場合でもトークン自体は有効かもしれない
			// TODO: エラーロギング
		}
	}

	return resp, nil
}

// RevokeToken は指定されたトークン (アクセスまたはリフレッシュ) を失効させます。
// RFC 7009 準拠。成功した場合は nil を返します。
// トークンが存在しない場合や既に失効している場合でもエラーとはしません。
func (s *TokenService) RevokeToken(ctx context.Context, tokenValue string, clientID domain.ClientID) error {
	token, err := s.tokenRepo.FindByValue(ctx, tokenValue)
	if err != nil {
		// トークンが見つからなくてもエラーではない
		if errors.Is(err, storage.ErrTokenNotFound) {
			return nil
		}
		// その他のリポジトリエラー
		// TODO: エラーロギング
		return NewOAuthError("server_error", "トークンの検索中にエラーが発生しました")
	}

	// クライアントIDが一致するか確認 (RFC 7009 Section 2.1)
	// トークンを発行したクライアントのみが失効できるようにする
	if token.ClientID != clientID {
		// エラーを返すか、単に何もしないか。RFCでは明確ではないが、エラーが親切か。
		return NewOAuthError("invalid_request", "トークンを発行したクライアントのみが失効できます")
	}

	// トークンを削除
	if err := s.tokenRepo.Delete(ctx, tokenValue); err != nil {
		// TODO: エラーロギング
		return NewOAuthError("server_error", "トークンの失効処理中にエラーが発生しました")
	}

	return nil
}

// --- ヘルパーメソッド ---

// authenticateClient はクライアントIDとシークレットでクライアントを認証します。
func (s *TokenService) authenticateClient(ctx context.Context, clientID domain.ClientID, clientSecret string) (domain.Client, error) {
	client, err := s.clientRepo.FindByID(ctx, clientID)
	if err != nil {
		return domain.Client{}, NewOAuthError("invalid_client", "指定されたクライアントIDは無効です")
	}

	// Public Client (シークレットなし)
	if client.Secret == "" {
		// Public Client がシークレットを提供してきた場合 (RFC 6749 Section 2.3.1)
		if clientSecret != "" {
			return domain.Client{}, NewOAuthError("invalid_client", "Public Client はクライアントシークレットを提供できません")
		}
		return client, nil // 認証成功
	}

	// Confidential Client (シークレットあり)
	if clientSecret == "" {
		return domain.Client{}, NewOAuthError("invalid_client", "Confidential Client はクライアントシークレットを提供する必要があります")
	}

	match, err := s.pwHasher.Compare(string(client.Secret), clientSecret)
	if err != nil {
		// ハッシュ比較エラー
		// TODO: エラーロギング
		return domain.Client{}, NewOAuthError("server_error", "クライアント認証中にエラーが発生しました")
	}
	if !match {
		return domain.Client{}, NewOAuthError("invalid_client", "クライアントシークレットが無効です")
	}

	return client, nil // 認証成功
}

// validateAuthorizationCode は認可コードを検証します。
func (s *TokenService) validateAuthorizationCode(ctx context.Context, codeValue string, clientID domain.ClientID, redirectURI string, now time.Time) (domain.AuthorizationCode, error) {
	if codeValue == "" {
		return domain.AuthorizationCode{}, NewOAuthError("invalid_grant", "認可コードが提供されていません")
	}

	authCode, err := s.codeRepo.FindByValue(ctx, codeValue)
	if err != nil {
		if errors.Is(err, storage.ErrCodeNotFound) {
			return domain.AuthorizationCode{}, NewOAuthError("invalid_grant", "無効な認可コードです")
		}
		// TODO: エラーロギング
		return domain.AuthorizationCode{}, NewOAuthError("server_error", "認可コードの検索中にエラーが発生しました")
	}

	// 有効期限チェック
	if authCode.IsExpired(now) {
		// 期限切れのコードは削除しておくのが望ましい
		_ = s.codeRepo.Delete(ctx, authCode.Value) // エラーは無視
		return domain.AuthorizationCode{}, NewOAuthError("invalid_grant", "認可コードの有効期限が切れています")
	}

	// クライアントIDの一致チェック
	if authCode.ClientID != clientID {
		return domain.AuthorizationCode{}, NewOAuthError("invalid_grant", "認可コードとクライアントIDが一致しません")
	}

	// リダイレクトURIの一致チェック (RFC 6749 Section 4.1.3)
	// 認可リクエスト時に redirect_uri が指定されていた場合のみ比較が必要
	// ここでは、コード発行時に保存したURIと、トークンリクエストで提供されたURIを比較する
	if authCode.RedirectURI != "" && authCode.RedirectURI != redirectURI {
		return domain.AuthorizationCode{}, NewOAuthError("invalid_grant", "リダイレクトURIが一致しません")
	}
	// 注意: トークンリクエストで redirect_uri が省略可能か、省略された場合にどうするかは仕様による

	// TODO: PKCE コードベリファイアの検証 (必要な場合)
	// if authCode.CodeChallenge != "" {
	// 	if !authCode.ValidatePKCE(codeVerifier) { // codeVerifier はリクエストから取得
	// 		return domain.AuthorizationCode{}, NewOAuthError("invalid_grant", "PKCEコードベリファイアが無効です")
	// 	}
	// }

	return authCode, nil
}

// authenticateUser はユーザー名とパスワードでユーザーを認証します。
func (s *TokenService) authenticateUser(ctx context.Context, username, password string) (domain.User, error) {
	if username == "" || password == "" {
		return domain.User{}, NewOAuthError("invalid_grant", "ユーザー名とパスワードは必須です")
	}
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		// ユーザーが見つからない
		return domain.User{}, NewOAuthError("invalid_grant", "ユーザー名またはパスワードが無効です")
	}
	match, err := s.pwHasher.Compare(user.HashedPassword, password)
	if err != nil {
		// ハッシュ比較エラー
		// TODO: エラーロギング
		return domain.User{}, NewOAuthError("server_error", "ユーザー認証中にエラーが発生しました")
	}
	if !match {
		// パスワード不一致
		return domain.User{}, NewOAuthError("invalid_grant", "ユーザー名またはパスワードが無効です")
	}
	return user, nil
}

// validateRefreshToken はリフレッシュトークンを検証します。
func (s *TokenService) validateRefreshToken(ctx context.Context, tokenValue string, clientID domain.ClientID, now time.Time) (domain.Token, error) {
	if tokenValue == "" {
		return domain.Token{}, NewOAuthError("invalid_grant", "リフレッシュトークンが提供されていません")
	}
	refreshToken, err := s.tokenRepo.FindByValue(ctx, tokenValue)
	if err != nil {
		if errors.Is(err, storage.ErrTokenNotFound) {
			return domain.Token{}, NewOAuthError("invalid_grant", "無効なリフレッシュトークンです")
		}
		// TODO: エラーロギング
		return domain.Token{}, NewOAuthError("server_error", "リフレッシュトークンの検索中にエラーが発生しました")
	}

	// 有効期限チェック
	if refreshToken.IsExpired(now) {
		// 期限切れのトークンは削除しておく
		_ = s.tokenRepo.Delete(ctx, refreshToken.Value) // エラーは無視
		return domain.Token{}, NewOAuthError("invalid_grant", "リフレッシュトークンの有効期限が切れています")
	}

	// クライアントIDの一致チェック
	if refreshToken.ClientID != clientID {
		return domain.Token{}, NewOAuthError("invalid_grant", "リフレッシュトークンとクライアントIDが一致しません")
	}

	// トークンタイプがリフレッシュトークンであることの確認 (もし区別している場合)
	// if refreshToken.Type != domain.TokenTypeRefresh { ... }

	return refreshToken, nil
}

// determineGrantedScopes は許可されるスコープを決定するヘルパー関数。
// clientScopes: クライアントに許可された全スコープ
// userScopes: ユーザーに許可された全スコープ (今回は未使用)
// requestedScopes: クライアントがリクエストしたスコープ
// 返り値: 実際に許可されるスコープ
func (s *TokenService) determineGrantedScopes(clientScopes []domain.Scope, userScopes []domain.Scope, requestedScopes []domain.Scope) []domain.Scope {
	allowedSet := make(map[domain.Scope]struct{})
	// クライアントに許可されたスコープをベースにする
	for _, s := range clientScopes {
		allowedSet[s] = struct{}{}
	}
	// TODO: ユーザー固有のスコープ制限があればここで考慮 (allowedSet から削除するなど)
	// if userScopes != nil { ... }

	var granted []domain.Scope
	if len(requestedScopes) == 0 {
		// リクエストがない場合は、クライアントに許可された全スコープを返すか、
		// あるいはサーバー定義のデフォルトスコープを返す。ここでは前者とする。
		// ただし、ユーザー固有の制限は考慮済み。
		for scope := range allowedSet {
			granted = append(granted, scope)
		}
		// 必要であればソートして順序を安定させる
		// sort.Slice(granted, func(i, j int) bool { return granted[i] < granted[j] })
	} else {
		// リクエストされたスコープのうち、許可されたものだけを返す
		seen := make(map[domain.Scope]struct{}) // 許可済みスコープの重複排除用
		for _, reqScope := range requestedScopes {
			if _, allowed := allowedSet[reqScope]; allowed {
				if _, alreadyAdded := seen[reqScope]; !alreadyAdded {
					granted = append(granted, reqScope)
					seen[reqScope] = struct{}{}
				}
			}
			// 許可されていないスコープがリクエストに含まれていてもエラーにはせず、単に無視する (RFC 6749 Section 3.3)
		}
	}
	// スコープが何も許可されなかった場合、空のスライスが返る

	// RFC 6749 Section 3.3: The authorization server SHOULD return the list of scopes granted
	// to the access token. -> granted スコープを返すのが推奨

	return granted
}
