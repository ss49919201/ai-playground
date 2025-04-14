package domain

import (
	"errors"
	"time"
	// "regexp" // Email validation if needed
)

// UserID はリソースオーナー（ユーザー）の一意な識別子です。
type UserID string

// User はリソースオーナーを表すエンティティです。
type User struct {
	ID             UserID
	Username       string    // ログインに使用されるユーザー名
	HashedPassword string    // ハッシュ化されて保存されるパスワード
	Email          string    // オプション: ユーザーのメールアドレス
	CreatedAt      time.Time // ユーザー作成日時
	// Scopes         []Scope   // オプション: ユーザーに紐づくデフォルトスコープや制限スコープ
}

// NewUser は新しい User エンティティを生成するファクトリ関数です。
// hashedPassword は既にハッシュ化されたパスワードを受け取ります。
// now は現在時刻を受け取り、CreatedAt フィールドに設定します (副作用の注入)。
// この関数は純粋関数として振る舞います。
func NewUser(id UserID, username, hashedPassword, email string, now time.Time) (User, error) {
	// --- バリデーション ---
	if id == "" {
		return User{}, errors.New("ユーザーIDは必須です")
	}
	if username == "" {
		return User{}, errors.New("ユーザー名は必須です")
	}
	if hashedPassword == "" {
		// ハッシュ化前のパスワードを受け取り、ここでハッシュ化する設計も考えられるが、
		// ファクトリ関数を純粋に保つため、ハッシュ化は呼び出し側で行う前提とする。
		return User{}, errors.New("ハッシュ化パスワードは必須です")
	}

	// オプション: Emailの形式バリデーション
	// if email != "" {
	// 	if !isValidEmail(email) {
	// 		return User{}, errors.New("無効なメールアドレス形式です")
	// 	}
	// }

	// --- エンティティ生成 ---
	return User{
		ID:             id,
		Username:       username,
		HashedPassword: hashedPassword,
		Email:          email,
		CreatedAt:      now,
	}, nil
}

/*
// isValidEmail は簡単なメールアドレス形式のチェックを行います。
// より厳密なチェックが必要な場合は、専用のライブラリを使用します。
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}
*/

// Note: パスワード検証 (Compare) は、セキュリティ上の理由と副作用（ハッシュ比較）を含むため、
// ドメインモデルのメソッドとしては実装せず、アプリケーションサービス層や
// PasswordHasher インターフェースを通じて行います。
