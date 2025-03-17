package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Message はチャットメッセージを表す値オブジェクト
type Message struct {
	ID        string    // メッセージの一意識別子
	Content   string    // メッセージ内容
	Sender    UserID    // 送信者ID
	Timestamp time.Time // 送信時刻
}

// NewMessage はメッセージオブジェクトを生成するファクトリ関数
func NewMessage(content string, sender UserID) (Message, error) {
	if content == "" {
		return Message{}, errors.New("メッセージ内容が空です")
	}
	return Message{
		ID:        uuid.New().String(),
		Content:   content,
		Sender:    sender,
		Timestamp: time.Now(),
	}, nil
}

// IsValid はメッセージが有効かどうかを検証するメソッド
func (m Message) IsValid() bool {
	return m.Content != "" && m.ID != ""
}

// FormattedTime はタイムスタンプのフォーマット済み文字列を返すメソッド
func (m Message) FormattedTime() string {
	return m.Timestamp.Format("15:04:05")
}

// Equal は2つのメッセージが等しいかどうかを判定するメソッド
func (m Message) Equal(other Message) bool {
	return m.ID == other.ID
}
