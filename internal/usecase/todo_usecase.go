package usecase

import (
	"context"
	"errors"
	"go_todo_app/internal/domain"
)

type TodoUseCase struct {
	repo TodoRepository
}

func NewTodoUseCase(repo TodoRepository) *TodoUseCase {
	return &TodoUseCase{repo: repo}
}

// CreateTodo TODOを作成するユースケース
func (uc *TodoUseCase) CreateTodo(ctx context.Context, title, description string) (*domain.Todo, error) {
	todo := &domain.Todo{
		Title:       title,
		Description: description,
		Completed:   false,
	}

	// ビジネスルールの検証
	if err := todo.Validate(); err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, todo); err != nil {
		return nil, err
	}

	return todo, nil
}

// GetAllTodos 全TODOを取得するユースケース
func (uc *TodoUseCase) GetAllTodos(ctx context.Context) ([]*domain.Todo, error) {
	return uc.repo.FindAll(ctx)
}

// GetTodoByID IDでTODOを取得するユースケース
func (uc *TodoUseCase) GetTodoByID(ctx context.Context, id int64) (*domain.Todo, error) {
	return uc.repo.FindByID(ctx, id)
}

// UpdateTodo TODOを更新するユースケース
func (uc *TodoUseCase) UpdateTodo(ctx context.Context, id int64, title, description string, completed bool) error {
	todo, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if todo == nil {
		return errors.New("TODOが見つかりません")
	}

	todo.Title = title
	todo.Description = description
	todo.Completed = completed

	if err := todo.Validate(); err != nil {
		return err
	}

	return uc.repo.Update(ctx, todo)
}

// ToggleTodo TODOの完了状態を切り替えるユースケース
func (uc *TodoUseCase) ToggleTodo(ctx context.Context, id int64) error {
	todo, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if todo == nil {
		return errors.New("TODOが見つかりません")
	}

	todo.Toggle()
	return uc.repo.Update(ctx, todo)
}

// DeleteTodo TODOを削除するユースケース
func (uc *TodoUseCase) DeleteTodo(ctx context.Context, id int64) error {
	return uc.repo.Delete(ctx, id)
}
