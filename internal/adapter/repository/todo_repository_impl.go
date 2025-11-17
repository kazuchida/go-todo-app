package repository

import (
	"context"
	"database/sql"
	"time"

	"go_todo_app/internal/domain"
)

type TodoRepositoryImpl struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) *TodoRepositoryImpl {
	return &TodoRepositoryImpl{db: db}
}

func (r *TodoRepositoryImpl) Create(ctx context.Context, todo *domain.Todo) error {
	query := `
        INSERT INTO todos (title, description, completed, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `
	now := time.Now()
	todo.CreatedAt = now
	todo.UpdatedAt = now

	return r.db.QueryRowContext(
		ctx, query,
		todo.Title, todo.Description, todo.Completed, todo.CreatedAt, todo.UpdatedAt,
	).Scan(&todo.ID)
}

func (r *TodoRepositoryImpl) FindByID(ctx context.Context, id int64) (*domain.Todo, error) {
	query := `
        SELECT id, title, description, completed, created_at, updated_at
        FROM todos WHERE id = $1
    `
	todo := &domain.Todo{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&todo.ID, &todo.Title, &todo.Description,
		&todo.Completed, &todo.CreatedAt, &todo.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (r *TodoRepositoryImpl) FindAll(ctx context.Context) ([]*domain.Todo, error) {
	query := `
        SELECT id, title, description, completed, created_at, updated_at
        FROM todos ORDER BY created_at DESC
    `
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*domain.Todo
	for rows.Next() {
		todo := &domain.Todo{}
		if err := rows.Scan(
			&todo.ID, &todo.Title, &todo.Description,
			&todo.Completed, &todo.CreatedAt, &todo.UpdatedAt,
		); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func (r *TodoRepositoryImpl) Update(ctx context.Context, todo *domain.Todo) error {
	query := `
        UPDATE todos
        SET title = $1, description = $2, completed = $3, updated_at = $4
        WHERE id = $5
    `
	todo.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(
		ctx, query,
		todo.Title, todo.Description, todo.Completed, todo.UpdatedAt, todo.ID,
	)
	return err
}

func (r *TodoRepositoryImpl) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM todos WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
