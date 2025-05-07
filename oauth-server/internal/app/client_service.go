package app

import (
	"context"
	"errors"
	"time"

	"github.com/ss49919201/ai-playground/go/oauth-server/internal/adapters/storage" // エラー型を参照するため
	"github.com/ss49919201/ai-playground/go/oauth-server/internal/domain"
	"github.com/ss49919201/ai-playground/go/oauth-server/internal/ports"
)

// ClientService はクライアントの登録や管理に関連するユースケースを処理します。
type ClientService struct {
	clientRepo   ports.ClientRepository
	idGenerator  ports.IDGenerator    // ClientID生成 (副作用)
	secretHasher ports.PasswordHasher // ClientSecretハッシュ化 (副作用)
	clock        ports.Clock          // 時刻取得 (副作用)
}

// NewClientService は ClientService の新しいインスタンスを生成します。
// 必要な依存関係 (リポジトリ、副作用インターフェースの実装) を引数として受け取ります。
func NewClientService(
	clientRepo ports.ClientRepository,
	idGenerator ports.IDGenerator,
	secretHasher ports.PasswordHasher,
	clock ports.Clock,
) *ClientService {
	return &ClientService{
		clientRepo:   clientRepo,
		idGenerator:  idGenerator,
		secretHasher: secretHasher,
		clock:        clock,
	}
}

// RegisterClientRequest はクライアント登録リクエストのパラメータです。
type RegisterClientRequest struct {
	Name         string   `json:"name"`          // クライアント名
	RedirectURIs []string `json:"redirect_uris"` // リダイレクトURIのリスト
	GrantTypes   []string `json:"grant_types"`   // 許可する認可フローのリスト
	Scopes       []string `json:"scopes"`        // 許可するスコープのリスト
}

// RegisterClientResponse はクライアント登録レスポンスのパラメータです。
// 生成された平文のクライアントシークレットを一度だけ含みます。
type RegisterClientResponse struct {
	ClientID     domain.ClientID `json:"client_id"`
	ClientSecret string          `json:"client_secret"` // 注意: このレスポンスでのみ返す
	Name         string          `json:"client_name"`
	RedirectURIs []string        `json:"redirect_uris"`
	GrantTypes   []string        `json:"grant_types"`
	Scopes       []string        `json:"scopes"`
	CreatedAt    time.Time       `json:"created_at"`
}

// RegisterClient は新しいクライアントを登録します。
// ClientIDとClientSecretを生成し、シークレットをハッシュ化して永続化します。
// 成功した場合、生成された情報（平文シークレットを含む）を返します。
func (s *ClientService) RegisterClient(ctx context.Context, req RegisterClientRequest) (RegisterClientResponse, error) {
	now := s.clock.Now()

	// 1. ClientID と ClientSecret の生成 (副作用)
	clientIDStr, err := s.idGenerator.Generate()
	if err != nil {
		// TODO: エラーロギング
		return RegisterClientResponse{}, errors.New("クライアントIDの生成に失敗しました")
	}
	clientID := domain.ClientID(clientIDStr)

	clientSecretPlain, err := s.idGenerator.GenerateSecret() // 平文のシークレット生成
	if err != nil {
		// TODO: エラーロギング
		return RegisterClientResponse{}, errors.New("クライアントシークレットの生成に失敗しました")
	}

	// 2. ClientSecret のハッシュ化 (副作用)
	hashedSecret, err := s.secretHasher.Hash(clientSecretPlain)
	if err != nil {
		// TODO: エラーロギング
		return RegisterClientResponse{}, errors.New("クライアントシークレットのハッシュ化に失敗しました")
	}

	// 3. GrantType と Scope のドメイン型への変換と検証
	grantTypes := make([]domain.GrantType, len(req.GrantTypes))
	for i, gtStr := range req.GrantTypes {
		gt := domain.GrantType(gtStr)
		// TODO: サポートされている GrantType かどうかの検証ロジックを追加
		// if !isValidGrantType(gt) {
		// 	return RegisterClientResponse{}, fmt.Errorf("無効なGrantTypeです: %s", gtStr)
		// }
		grantTypes[i] = gt
	}

	scopes := make([]domain.Scope, len(req.Scopes))
	for i, scStr := range req.Scopes {
		// ValidateScope で形式チェックは行われるが、ここでは存在チェックなどを行う場合がある
		// TODO: サーバーで定義されているスコープかどうかの検証ロジックを追加
		scopes[i] = domain.Scope(scStr)
	}
	// ドメインの ValidateScope を使って文字列からパース＆バリデーションする方が良いかもしれない
	// parsedScopes, err := domain.ValidateScope(strings.Join(req.Scopes, " "))
	// if err != nil {
	// 	return RegisterClientResponse{}, fmt.Errorf("無効なスコープ形式です: %w", err)
	// }

	// 4. ドメインオブジェクト (Client) の生成
	client, err := domain.NewClient(clientID, domain.ClientSecret(hashedSecret), req.Name, req.RedirectURIs, grantTypes, scopes, now)
	if err != nil {
		// ドメインレベルのバリデーションエラー
		return RegisterClientResponse{}, err
	}

	// 5. リポジトリに保存 (副作用)
	if err := s.clientRepo.Save(ctx, client); err != nil {
		// TODO: エラーロギング
		// 重複エラーなどをハンドリングする必要がある場合がある
		return RegisterClientResponse{}, errors.New("クライアント情報の保存に失敗しました")
	}

	// 6. レスポンス生成 (平文のシークレットを含む)
	resp := RegisterClientResponse{
		ClientID:     client.ID,
		ClientSecret: clientSecretPlain, // 平文シークレットを返す
		Name:         client.Name,
		RedirectURIs: client.RedirectURIs,
		GrantTypes:   req.GrantTypes, // 元のリクエストの文字列スライスを返す
		Scopes:       req.Scopes,     // 元のリクエストの文字列スライスを返す
		CreatedAt:    client.CreatedAt,
	}

	return resp, nil
}

// GetClientResponse はクライアント情報取得レスポンスのパラメータです。
// セキュリティのため、クライアントシークレットは含みません。
type GetClientResponse struct {
	ClientID     domain.ClientID `json:"client_id"`
	Name         string          `json:"client_name"`
	RedirectURIs []string        `json:"redirect_uris"`
	GrantTypes   []string        `json:"grant_types"`
	Scopes       []string        `json:"scopes"`
	CreatedAt    time.Time       `json:"created_at"`
}

// GetClient は指定されたクライアントIDのクライアント情報を取得します。
// 取得した情報からクライアントシークレットを除外して返します。
func (s *ClientService) GetClient(ctx context.Context, clientID domain.ClientID) (GetClientResponse, error) {
	client, err := s.clientRepo.FindByID(ctx, clientID)
	if err != nil {
		if errors.Is(err, storage.ErrClientNotFound) {
			// Not Found エラーを返す (HTTP層で404に変換されることを想定)
			return GetClientResponse{}, err // storageのエラーをそのまま返すか、app層のエラーにラップするか
		}
		// その他のリポジトリエラー
		// TODO: エラーロギング
		return GetClientResponse{}, errors.New("クライアント情報の取得に失敗しました")
	}

	// レスポンス用に変換 (シークレットを除外)
	grantTypesStr := make([]string, len(client.GrantTypes))
	for i, gt := range client.GrantTypes {
		grantTypesStr[i] = string(gt)
	}
	scopesStr := make([]string, len(client.Scopes))
	for i, sc := range client.Scopes {
		scopesStr[i] = string(sc)
	}

	resp := GetClientResponse{
		ClientID:     client.ID,
		Name:         client.Name,
		RedirectURIs: client.RedirectURIs,
		GrantTypes:   grantTypesStr,
		Scopes:       scopesStr,
		CreatedAt:    client.CreatedAt,
	}

	return resp, nil
}

// AuthenticateClient は提供されたクライアントIDとシークレットを使用してクライアントを認証します。
// 認証に成功した場合は Client エンティティを、失敗した場合はエラーを返します。
// このメソッドは主にトークンエンドポイントで使用されます。
func (s *ClientService) AuthenticateClient(ctx context.Context, clientID domain.ClientID, clientSecret string) (domain.Client, error) {
	client, err := s.clientRepo.FindByID(ctx, clientID)
	if err != nil {
		// クライアントが見つからない場合は認証失敗
		return domain.Client{}, errors.New("クライアント認証に失敗しました (invalid_client)")
	}

	// public クライアント (シークレットが設定されていない) の場合、シークレット検証は不要
	// ただし、シークレットが必要なフロー (例: authorization_code) でシークレットが提供されなかった場合はエラーとする必要がある
	// このロジックは呼び出し側 (TokenService) で行う方が適切かもしれない

	if client.Secret == "" {
		// public クライアントの場合、シークレットが提供されていたらエラー？ 仕様による
		if clientSecret != "" {
			return domain.Client{}, errors.New("クライアント認証に失敗しました (public clientにシークレットが提供されました)")
		}
		// public クライアントとして認証成功
		return client, nil
	}

	// confidential クライアントの場合、シークレットを比較
	match, err := s.secretHasher.Compare(string(client.Secret), clientSecret)
	if err != nil {
		// ハッシュ比較エラー (ログ出力推奨)
		// TODO: エラーロギング
		return domain.Client{}, errors.New("クライアント認証中にエラーが発生しました")
	}
	if !match {
		// シークレット不一致
		return domain.Client{}, errors.New("クライアント認証に失敗しました (invalid_client)")
	}

	// 認証成功
	return client, nil
}
