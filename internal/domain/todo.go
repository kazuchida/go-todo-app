package domain

import (
	"errors"
	"time"
)

// Todo エンティティ(ビジネスルールの中心)
type Todo struct {
	ID          int64
	Title       string
	Description string
	Completed   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ビジネスルール: バリデーション
func (t *Todo) Validate() error {
	if t.Title == "" {
		return errors.New("タイトルは必須です")
	}
	if len(t.Title) > 100 {
		return errors.New("タイトルは100文字以内です")
	}
	return nil
}

// ビジネスルール: 完了状態の切り替え
func (t *Todo) Toggle() {
	t.Completed = !t.Completed
	t.UpdatedAt = time.Now()
}
