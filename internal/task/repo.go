package task

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("Task not found")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, task *Task) error {
	query := `
		INSERT INTO tasks (title, description, status, due_date, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Status,
		task.DueDate,
		task.CreatedAt,
		task.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert task: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	task.ID = id

	return nil
}

func (r *Repository) GetById(ctx context.Context, id int64) (*Task, error) {

	query := `SELECT id, title, description, status, due_date, created_at, updated_at FROM tasks WHERE id = ?`

	row := r.db.QueryRowContext(ctx, query, id)

	var task Task

	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.DueDate, &task.CreatedAt, &task.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get task by id: %w", err)
	}

	return &task, nil
}

func (r *Repository) GetAll(ctx context.Context) ([]Task, error) {

	query := `
		SELECT id, title, description, status, due_date, created_at, updated_at 
		FROM tasks 
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}

	defer rows.Close()

	tasks := make([]Task, 0)

	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.DueDate, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task row: %w", err)
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows itteration error: %w", err)
	}

	return tasks, nil
}

func (r *Repository) update(ctx context.Context, task *Task)  error {

	query := `UPDATE tasks
	SET title = ?, description = ? , status = ?, due_date = ? , updated_at = ?
	WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx,query, task.Title, task.Description, task.Status, task.DueDate, task.UpdatedAt, task.ID)

	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}


	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM tasks WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)

	if err != nil {
		return fmt.Errorf("failed to delete task %w", err)
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}