package domain

import (
	"errors"
	"time"
)

// AuthorizationCode は認可コードフローで使用される一時的なコードを表す値オブジェクトです。
// 一度使用されると無効になります。イミュータブルとして扱います。
type AuthorizationCode struct {
	Value       string    // 認可コード自体の値 (ランダム文字列)
	ClientID    ClientID  // このコードを発行されたクライアントのID
	UserID      UserID    // このコードに関連付けられたユーザーのID
	RedirectURI string    // 認可リクエストで使用されたリダイレクトURI
	Scopes      []Scope   // 認可リクエストで許可されたスコープ
	ExpiresAt   time.Time // 認可コードの有効期限 (通常は短い、例: 10分)
	IssuedAt    time.Time // 認可コードの発行日時
	// --- PKCE (Proof Key for Code Exchange) 関連フィールド (オプション) ---
	// CodeChallenge       string // PKCE コードチャレンジ (S256ハッシュなど)
	// CodeChallengeMethod string // PKCE チャレンジメソッド ("S256" または "plain")
}

// NewAuthorizationCode は新しい AuthorizationCode 値オブジェクトを生成するファクトリ関数です。
// value は認可コード文字列、clientID, userID, redirectURI, scopes は関連情報、
// issuedAt, expiresAt は発行日時と有効期限です。
// この関数は純粋関数として振る舞います。
func NewAuthorizationCode(value string, clientID ClientID, userID UserID, redirectURI string, scopes []Scope, issuedAt, expiresAt time.Time /*, codeChallenge, codeChallengeMethod string*/) (AuthorizationCode, error) {
	// --- バリデーション ---
	if value == "" {
		return AuthorizationCode{}, errors.New("認可コードの値は必須です")
	}
	if clientID == "" {
		return AuthorizationCode{}, errors.New("クライアントIDは必須です")
	}
	if userID == "" {
		// 認可コードは必ずユーザーに紐づくはず
		return AuthorizationCode{}, errors.New("ユーザーIDは必須です")
	}
	if redirectURI == "" {
		// リダイレクトURIも必須 (トークンリクエスト時の検証に使用)
		return AuthorizationCode{}, errors.New("リダイレクトURIは必須です")
	}
	if expiresAt.Before(issuedAt) {
		return AuthorizationCode{}, errors.New("認可コードの有効期限が発行日時より前です")
	}
	// PKCE関連のバリデーション (オプション)
	// if (codeChallenge != "" && codeChallengeMethod == "") || (codeChallenge == "" && codeChallengeMethod != "") {
	// 	return AuthorizationCode{}, errors.New("PKCEを使用する場合、コードチャレンジとメソッドの両方が必要です")
	// }
	// if codeChallengeMethod != "" && codeChallengeMethod != "S256" && codeChallengeMethod != "plain" {
	// 	return AuthorizationCode{}, errors.New("無効なPKCEコードチャレンジメソッドです")
	// }

	// --- 値オブジェクト生成 ---
	scopesCopy := make([]Scope, len(scopes))
	copy(scopesCopy, scopes)

	return AuthorizationCode{
		Value:       value,
		ClientID:    clientID,
		UserID:      userID,
		RedirectURI: redirectURI,
		Scopes:      scopesCopy,
		ExpiresAt:   expiresAt,
		IssuedAt:    issuedAt,
		// CodeChallenge:       codeChallenge,
		// CodeChallengeMethod: codeChallengeMethod,
	}, nil
}

// IsExpired は指定された時刻 (now) において認可コードが有効期限切れかどうかを返します。
// このメソッドは純粋関数です。
func (c AuthorizationCode) IsExpired(now time.Time) bool {
	if c.ExpiresAt.IsZero() {
		return false // 通常は有効期限必須
	}
	return !now.Before(c.ExpiresAt)
}

/*
// ValidatePKCE は提供されたコードベリファイアが、保存されたコードチャレンジと一致するか検証します。
// このメソッド自体は純粋ではありません（ハッシュ計算を含むため）。
// 通常、この検証はアプリケーションサービス層で行われます。
func (c AuthorizationCode) ValidatePKCE(codeVerifier string) (bool, error) {
	if c.CodeChallenge == "" || c.CodeChallengeMethod == "" {
		// PKCEが使用されていない場合は検証不要 (あるいはエラー？)
		return true, nil // またはエラーを返すか、仕様による
	}

	switch c.CodeChallengeMethod {
	case "plain":
		return c.CodeChallenge == codeVerifier, nil
	case "S256":
		// codeVerifier を SHA256 でハッシュ化し、Base64 URLエンコードした結果と比較
		h := sha256.Sum256([]byte(codeVerifier))
		encoded := base64.RawURLEncoding.EncodeToString(h[:]) // RawURLEncoding を使用
		return c.CodeChallenge == encoded, nil
	default:
		return false, errors.New("不明なPKCEコードチャレンジメソッドです: " + c.CodeChallengeMethod)
	}
}
*/
