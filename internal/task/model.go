package task

import "time"

type Status string

const (
	StatusPending Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusCompleted Status = "completed"
)

type Task struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description *string     `json:"description"`
	Status      Status     `json:"status"`
	DueDate     *time.Time `json:"due_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateRequest struct {
	Title string `json:"title"`
	Description *string `json:"description"`
	DueDate *time.Time `json:"due_date"`
}

type UpdateRequest struct {
	Title string `json:"title"`
	Description *string `json:"description"`
	DueDate *time.Time `json:"due_date"`
	Status      Status     `json:"status"`
}