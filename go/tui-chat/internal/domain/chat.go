package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// 定義済みエラー
var (
	ErrSessionFull      = errors.New("チャットセッションは最大2人までです")
	ErrUserNotInSession = errors.New("ユーザーはセッションに参加していません")
	ErrDuplicateUser    = errors.New("ユーザーは既にセッションに参加しています")
)

// ChatSession はチャットセッションを表す集約ルート
type ChatSession struct {
	ID        string
	Users     []User
	Messages  []Message
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewChatSession はチャットセッションを生成するファクトリ関数
func NewChatSession(initialUser User) ChatSession {
	return ChatSession{
		ID:        uuid.New().String(),
		Users:     []User{initialUser},
		Messages:  []Message{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// AddUser はチャットセッションにユーザーを追加するメソッド
func (c *ChatSession) AddUser(user User) error {
	if len(c.Users) >= 2 {
		return ErrSessionFull
	}

	// 既に参加しているユーザーかチェック
	for _, u := range c.Users {
		if u.ID == user.ID {
			return ErrDuplicateUser
		}
	}

	c.Users = append(c.Users, user)
	c.UpdatedAt = time.Now()
	return nil
}

// RemoveUser はチャットセッションからユーザーを削除するメソッド
func (c *ChatSession) RemoveUser(userID UserID) error {
	for i, user := range c.Users {
		if user.ID == userID {
			// ユーザーを削除
			c.Users = append(c.Users[:i], c.Users[i+1:]...)
			c.UpdatedAt = time.Now()
			return nil
		}
	}
	return ErrUserNotInSession
}

// AddMessage はチャットセッションにメッセージを追加するメソッド
func (c *ChatSession) AddMessage(msg Message) error {
	// 送信者がセッションに参加しているか確認
	senderExists := false
	for _, user := range c.Users {
		if user.ID == msg.Sender {
			senderExists = true
			break
		}
	}
	if !senderExists {
		return ErrUserNotInSession
	}

	c.Messages = append(c.Messages, msg)
	c.UpdatedAt = time.Now()
	return nil
}

// GetMessages はチャットセッションのメッセージを取得するメソッド
func (c *ChatSession) GetMessages() []Message {
	// メッセージのコピーを返す
	messages := make([]Message, len(c.Messages))
	copy(messages, c.Messages)
	return messages
}

// GetUsers はチャットセッションのユーザーを取得するメソッド
func (c *ChatSession) GetUsers() []User {
	// ユーザーのコピーを返す
	users := make([]User, len(c.Users))
	copy(users, c.Users)
	return users
}

// IsFull はチャットセッションが満員かどうかを判定するメソッド
func (c *ChatSession) IsFull() bool {
	return len(c.Users) >= 2
}

// GetUserByID はIDでユーザーを検索するメソッド
func (c *ChatSession) GetUserByID(id UserID) (User, bool) {
	for _, user := range c.Users {
		if user.ID == id {
			return user, true
		}
	}
	return User{}, false
}
