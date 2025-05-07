package domain

import (
	"errors"
	"time"
)

// ClientID はクライアントの一意な識別子です。
type ClientID string

// ClientSecret はクライアントのシークレットです。
// 永続化する際はハッシュ化された値を想定します。
type ClientSecret string

// GrantType はクライアントに許可される認可フローの種類を表します。
type GrantType string

// 定義済みの GrantType
const (
	GrantTypeAuthorizationCode GrantType = "authorization_code"
	GrantTypeImplicit          GrantType = "implicit" // オプション (非推奨)
	GrantTypeClientCredentials GrantType = "client_credentials"
	GrantTypePassword          GrantType = "password" // オプション (非推奨)
	GrantTypeRefreshToken      GrantType = "refresh_token"
)

// Client は OAuth 2.0 クライアントアプリケーションを表すエンティティです。
type Client struct {
	ID           ClientID
	Secret       ClientSecret // 永続化層でハッシュ化された値が格納されることを想定
	Name         string
	RedirectURIs []string    // 認可コード/インプリシットフローで使用されるリダイレクト先URI
	GrantTypes   []GrantType // このクライアントが使用を許可されている認可フロー
	Scopes       []Scope     // このクライアントが要求を許可されているスコープ
	CreatedAt    time.Time   // クライアント作成日時
}

// NewClient は新しい Client エンティティを生成するファクトリ関数です。
// hashedSecret は既にハッシュ化されたクライアントシークレットを受け取ります。
// now は現在時刻を受け取り、CreatedAt フィールドに設定します (副作用の注入)。
// この関数は純粋関数として振る舞います (引数が同じなら常に同じ結果を返す)。
func NewClient(id ClientID, hashedSecret ClientSecret, name string, redirectURIs []string, grantTypes []GrantType, scopes []Scope, now time.Time) (Client, error) {
	// --- バリデーション ---
	if id == "" {
		return Client{}, errors.New("クライアントIDは必須です")
	}
	if name == "" {
		return Client{}, errors.New("クライアント名は必須です")
	}
	// hashedSecret のバリデーションはここでは行わない (空を許可する場合もあるため)

	// RedirectURIs のバリデーション (少なくとも1つ必要か、形式は正しいかなど)
	// GrantTypes のバリデーション (サポートされているタイプか)
	// Scopes のバリデーション (空でも良いかなど)

	// GrantType に応じた RedirectURIs の要件チェック
	requiresRedirectURI := false
	for _, gt := range grantTypes {
		if gt == GrantTypeAuthorizationCode || gt == GrantTypeImplicit {
			requiresRedirectURI = true
			break
		}
	}
	if requiresRedirectURI && len(redirectURIs) == 0 {
		return Client{}, errors.New("認可コードフローまたはインプリシットフローを使用する場合、リダイレクトURIが少なくとも1つ必要です")
	}
	// TODO: redirectURIs の形式 (有効なURIか、フラグメントを含まないかなど) を検証

	// --- エンティティ生成 ---
	// コピーを作成してイミュータビリティを確保
	urisCopy := make([]string, len(redirectURIs))
	copy(urisCopy, redirectURIs)
	grantsCopy := make([]GrantType, len(grantTypes))
	copy(grantsCopy, grantTypes)
	scopesCopy := make([]Scope, len(scopes))
	copy(scopesCopy, scopes)

	return Client{
		ID:           id,
		Secret:       hashedSecret,
		Name:         name,
		RedirectURIs: urisCopy,
		GrantTypes:   grantsCopy,
		Scopes:       scopesCopy,
		CreatedAt:    now,
	}, nil
}

// ValidateRedirectURI は指定されたURIがクライアントに登録されたリダイレクトURIのいずれかと
// 一致するかどうかを検証します。
// このメソッドは純粋関数です。
func (c Client) ValidateRedirectURI(uri string) bool {
	if uri == "" {
		return false // URIが空の場合は常に無効
	}
	for _, registeredURI := range c.RedirectURIs {
		if registeredURI == uri {
			return true
		}
	}
	return false
}

// HasGrantType は指定された認可フロー (GrantType) が
// このクライアントに許可されているかどうかを検証します。
// このメソッドは純粋関数です。
func (c Client) HasGrantType(grantType GrantType) bool {
	for _, allowedType := range c.GrantTypes {
		if allowedType == grantType {
			return true
		}
	}
	// リフレッシュトークンは独立したGrantTypeとして扱われることが多いが、
	// 他のフローの結果として使用されるため、明示的に許可リストに含まれていなくても
	// 許可されているとみなす場合もある。ここでは明示的な許可が必要とする。
	return false
}

// ValidateScope は要求されたスコープのリストが、
// このクライアントに許可されたスコープの範囲内であるかどうかを検証します。
// 要求されたスコープが空の場合は常に true を返します。
// 要求されたスコープの中に、クライアントに許可されていないスコープが1つでも含まれていれば false を返します。
// このメソッドは純粋関数です。
func (c Client) ValidateScope(requestedScopes []Scope) bool {
	if len(requestedScopes) == 0 {
		return true // スコープが要求されていない場合は常に有効
	}

	// クライアントに許可されたスコープをセットに変換して高速検索
	allowedScopeSet := make(map[Scope]struct{}, len(c.Scopes))
	for _, s := range c.Scopes {
		allowedScopeSet[s] = struct{}{}
	}

	// 要求された各スコープが許可されているかチェック
	for _, reqScope := range requestedScopes {
		if _, ok := allowedScopeSet[reqScope]; !ok {
			return false // 許可されていないスコープが含まれている
		}
	}
	return true
}

// ValidateSecret は提供された平文のシークレットが、クライアントの（ハッシュ化された）シークレットと
// 一致するかどうかを検証します。実際の比較は PasswordHasher インターフェースが行います。
// このメソッド自体は純粋ではありません（PasswordHasherに依存するため）。
// 通常、この検証はアプリケーションサービス層で行われます。
// func (c Client) ValidateSecret(plainSecret string, hasher ports.PasswordHasher) (bool, error) {
// 	if c.Secret == "" {
// 		// シークレットが設定されていないクライアント (public client)
// 		return false, errors.New("クライアントシークレットが設定されていません")
// 	}
// 	return hasher.Compare(string(c.Secret), plainSecret)
// }
