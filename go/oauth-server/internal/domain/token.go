package domain

import (
	"errors"
	"time"
)

// TokenType は発行されるトークンの種類を示します。
// OAuth 2.0 では通常 "Bearer" が使用されます。
type TokenType string

const (
	TokenTypeBearer TokenType = "Bearer"
)

// Token はアクセストークンまたはリフレッシュトークンを表す値オブジェクトです。
// イミュータブル（不変）として扱います。
type Token struct {
	Value     string    // トークン自体の値 (JWT形式またはランダム文字列)
	Type      TokenType // トークンタイプ (例: "Bearer")
	ClientID  ClientID  // このトークンを発行されたクライアントのID
	UserID    UserID    // このトークンに関連付けられたユーザーのID (クライアントクレデンシャル フローの場合は空になることがある)
	Scopes    []Scope   // このトークンに許可されたスコープ
	ExpiresAt time.Time // トークンの有効期限
	IssuedAt  time.Time // トークンの発行日時
}

// NewToken は新しい Token 値オブジェクトを生成するファクトリ関数です。
// value はトークン文字列、tokenType はトークンタイプ (通常 TokenTypeBearer)、
// clientID, userID, scopes はトークンに関連付けられる情報、
// issuedAt, expiresAt は発行日時と有効期限です。
// この関数は純粋関数として振る舞います。
func NewToken(value string, tokenType TokenType, clientID ClientID, userID UserID, scopes []Scope, issuedAt, expiresAt time.Time) (Token, error) {
	// --- バリデーション ---
	if value == "" {
		return Token{}, errors.New("トークン値は必須です")
	}
	if clientID == "" {
		// クライアントクレデンシャル フローでも ClientID は必須
		return Token{}, errors.New("クライアントIDは必須です")
	}
	// UserID は空を許可する (クライアントクレデンシャル フローの場合)
	if expiresAt.Before(issuedAt) {
		return Token{}, errors.New("トークンの有効期限が発行日時より前です")
	}
	if tokenType == "" {
		// デフォルト値を設定するか、エラーとするか
		tokenType = TokenTypeBearer
	} else if tokenType != TokenTypeBearer {
		// 現状 Bearer のみを想定
		return Token{}, errors.New("無効なトークンタイプです: " + string(tokenType))
	}

	// --- 値オブジェクト生成 ---
	// スコープのスライスをコピーして不変性を保証
	scopesCopy := make([]Scope, len(scopes))
	copy(scopesCopy, scopes)

	return Token{
		Value:     value,
		Type:      tokenType,
		ClientID:  clientID,
		UserID:    userID,
		Scopes:    scopesCopy,
		ExpiresAt: expiresAt,
		IssuedAt:  issuedAt,
	}, nil
}

// IsExpired は指定された時刻 (now) においてトークンが有効期限切れかどうかを返します。
// このメソッドは純粋関数です。
func (t Token) IsExpired(now time.Time) bool {
	// ExpiresAt がゼロ値の場合は有効期限なしとみなすか？ -> OAuthでは通常有効期限必須
	if t.ExpiresAt.IsZero() {
		return false // 有効期限が設定されていない場合は常に有効？ 仕様によるが、通常はエラー
	}
	// ExpiresAt と now が同じ時刻の場合も期限切れとみなす（After は含まないため）
	return !now.Before(t.ExpiresAt)
}

// HasScope はトークンが必要なスコープを含んでいるかどうかを返します。
// このメソッドは純粋関数です。
func (t Token) HasScope(requiredScope Scope) bool {
	for _, scope := range t.Scopes {
		if scope == requiredScope {
			return true
		}
	}
	return false
}

// HasAllScopes はトークンが必要なスコープすべてを含んでいるかどうかを返します。
// requiredScopes が空の場合は true を返します。
// このメソッドは純粋関数です。
func (t Token) HasAllScopes(requiredScopes []Scope) bool {
	if len(requiredScopes) == 0 {
		return true
	}

	tokenScopesSet := make(map[Scope]struct{}, len(t.Scopes))
	for _, s := range t.Scopes {
		tokenScopesSet[s] = struct{}{}
	}

	for _, reqScope := range requiredScopes {
		if _, ok := tokenScopesSet[reqScope]; !ok {
			return false // 必要なスコープが1つでも欠けていれば false
		}
	}
	return true
}
