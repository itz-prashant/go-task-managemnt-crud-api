package task

import (
	"context"
	"errors"
	"strings"
	"time"
)

var (
	ErrEmptyTitle    = errors.New("Title can not be empty")
	ErrInvalidId     = errors.New("invalid id: must be greater than 0")
	ErrInvalidStatus = errors.New("invalid status: must be pending, in_progress, or completed")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateTask(ctx context.Context, req CreateRequest) (*Task, error) {
	if strings.TrimSpace(req.Title) == "" {
		return nil, ErrEmptyTitle
	}

	currentTime := time.Now().UTC()

	task := &Task{
		Title:       strings.TrimSpace(req.Title),
		Description: req.Description,
		Status:      StatusPending,
		DueDate:     req.DueDate,
		CreatedAt:   currentTime,
		UpdatedAt:   currentTime,
	}

	err := s.repo.Create(ctx, task)

	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *Service) GetTask(ctx context.Context, id int64) (*Task, error) {
	if id <= 0 {
		return nil, ErrInvalidId
	}
	return s.repo.GetById(ctx, id)
}

func (s *Service) ListTask(ctx context.Context) ([]Task, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) UpdateTask(ctx context.Context, id int64, req UpdateRequest) (*Task, error) {
	if id <= 0 {
		return nil, ErrInvalidId
	}

	if strings.TrimSpace(req.Title) == "" {
		return nil, ErrEmptyTitle
	}

	if req.Status != StatusPending && req.Status != StatusInProgress && req.Status != StatusCompleted {
		return nil, ErrInvalidStatus
	}

	task, err := s.repo.GetById(ctx, id)

	if err != nil {
		return nil, err
	}

	task.Title = req.Title
	task.Status = req.Status
	task.Description = req.Description
	task.DueDate = req.DueDate
	task.UpdatedAt = time.Now().UTC()

	err = s.repo.update(ctx, task)

	return task, nil
}

func (s *Service) DeleteTask (ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrInvalidId
	}

	return s.repo.Delete(ctx, id)
}