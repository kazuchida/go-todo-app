package usecase

import (
	"context"
	"go_todo_app/internal/domain"
)

// TodoRepository リポジトリのインターフェース(依存性逆転)
type TodoRepository interface {
	Create(ctx context.Context, todo *domain.Todo) error
	FindByID(ctx context.Context, id int64) (*domain.Todo, error)
	FindAll(ctx context.Context) ([]*domain.Todo, error)
	Update(ctx context.Context, todo *domain.Todo) error
	Delete(ctx context.Context, id int64) error
}
