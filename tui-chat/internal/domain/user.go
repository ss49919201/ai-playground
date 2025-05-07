package domain

import (
	"errors"

	"github.com/google/uuid"
)

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

// String はUserIDの文字列表現を返すメソッド
func (id UserID) String() string {
	return string(id)
}

// Deactivate はユーザーを非アクティブ状態にするメソッド
func (u *User) Deactivate() {
	u.IsActive = false
}

// Activate はユーザーをアクティブ状態にするメソッド
func (u *User) Activate() {
	u.IsActive = true
}
