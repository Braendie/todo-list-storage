package storage

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Braendie/todo-list-storage/internal/models"
	"github.com/Braendie/todo-list-storage/internal/storage"
)

type StorageService struct {
	log     *slog.Logger
	storage TaskStorage
}

type TaskStorage interface {
	TaskCreater
	TaskGetter
	TaskDeleter
	TaskUpdater
}

type TaskCreater interface {
	CreateTask(ctx context.Context, title string, description string) (int64, error)
}

type TaskGetter interface {
	GetTasks(ctx context.Context) ([]models.Task, error)
}

type TaskDeleter interface {
	DeleteTask(ctx context.Context, id int64) error
}

type TaskUpdater interface {
	UpdateTask(ctx context.Context, id int64) error
}

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

func New(log *slog.Logger, storage TaskStorage) *StorageService {
	return &StorageService{
		log:     log,
		storage: storage,
	}
}

func (s *StorageService) CreateTask(ctx context.Context, title string, description string) (int64, error) {
	const op = "service.StorageService.CreateTask"

	s.log.With(
		slog.String("op", op),
	)

	id, err := s.storage.CreateTask(ctx, title, description)
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			s.log.Warn("Task already exists", slog.String("op", op), slog.String("title", title))
			return 0, ErrAlreadyExists
		}

		s.log.Error("Failed to create task", slog.String("op", op), slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *StorageService) GetTasks(ctx context.Context) ([]models.Task, error) {
	const op = "service.StorageService.GetTasks"

	s.log.With(
		slog.String("op", op),
	)

	tasks, err := s.storage.GetTasks(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			s.log.Info("No tasks found", slog.String("op", op))
			return nil, ErrNotFound
		}

		s.log.Error("Failed to get tasks", slog.String("op", op), slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return tasks, nil
}

func (s *StorageService) DeleteTask(ctx context.Context, id int64) error {
	const op = "service.StorageService.DeleteTask"

	s.log.With(
		slog.String("op", op),
		slog.Int64("id", id),
	)

	err := s.storage.DeleteTask(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			s.log.Warn("Task not found", slog.String("op", op), slog.Int64("id", id))
			return ErrNotFound
		}

		s.log.Error("Failed to delete task", slog.String("op", op), slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *StorageService) UpdateTask(ctx context.Context, id int64) error {
	const op = "service.StorageService.UpdateTask"

	s.log.With(
		slog.String("op", op),
		slog.Int64("id", id),
	)

	err := s.storage.UpdateTask(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			s.log.Warn("Task not found for update", slog.String("op", op), slog.Int64("id", id))
			return ErrNotFound
		}

		s.log.Error("Failed to update task", slog.String("op", op), slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
