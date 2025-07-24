package storage

import (
	"context"
	"errors"

	storagev1 "github.com/Braendie/todo-list-protos/gen/go/storage"
	"github.com/Braendie/todo-list-storage/internal/models"
	storageService "github.com/Braendie/todo-list-storage/internal/services/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Storage interface {
	CreateTask(ctx context.Context, title string, description string) (int64, error)
	GetTasks(ctx context.Context) ([]models.Task, error)
	DeleteTask(ctx context.Context, id int64) error
	UpdateTask(ctx context.Context, id int64) error
}

type serverAPI struct {
	storagev1.UnimplementedStorageServer
	storage Storage
}

func Register(gRPC *grpc.Server, storageService Storage) {
	storagev1.RegisterStorageServer(gRPC, &serverAPI{storage: storageService})
}

func (s *serverAPI) Create(ctx context.Context, req *storagev1.CreateRequest) (*storagev1.CreateResponse, error) {
	if err := validateCreate(req); err != nil {
		return nil, err
	}

	taskID, err := s.storage.CreateTask(ctx, req.GetTitle(), req.GetDescription())
	if err != nil {
		if errors.Is(err, storageService.ErrAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "task already exists")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &storagev1.CreateResponse{
		TaskId: taskID,
	}, nil
}

func (s *serverAPI) Delete(ctx context.Context, req *storagev1.DeleteRequest) (*emptypb.Empty, error) {
	if err := validateDelete(req); err != nil {
		return nil, err
	}

	err := s.storage.DeleteTask(ctx, req.TaskId)
	if err != nil {
		if errors.Is(err, storageService.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "task not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return nil, nil
}

func (s *serverAPI) Done(ctx context.Context, req *storagev1.DoneRequest) (*emptypb.Empty, error) {
	if err := validateDone(req); err != nil {
		return nil, err
	}

	err := s.storage.UpdateTask(ctx, req.TaskId)
	if err != nil {
		if errors.Is(err, storageService.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "task not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return nil, nil
}

func (s *serverAPI) List(ctx context.Context, _ *emptypb.Empty) (*storagev1.ListResponse, error) {
	tasks, err := s.storage.GetTasks(ctx)
	if err != nil {
		if errors.Is(err, storageService.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "tasks not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	resTasks := make([]*storagev1.Task, 0, len(tasks))
	for _, val := range tasks {
		resTasks = append(resTasks, &storagev1.Task{
			TaskId:      val.ID,
			Title:       val.Title,
			Description: val.Description,
			Done:        val.Done,
		})
	}

	return &storagev1.ListResponse{
		Tasks: resTasks,
	}, nil
}

func validateCreate(req *storagev1.CreateRequest) error {
	if req.GetTitle() == "" {
		return status.Error(codes.InvalidArgument, "title is required")
	}

	return nil
}

func validateDelete(req *storagev1.DeleteRequest) error {
	if req.TaskId < 1 {
		return status.Error(codes.InvalidArgument, "id must be greater than 0")
	}

	return nil
}

func validateDone(req *storagev1.DoneRequest) error {
	if req.TaskId < 1 {
		return status.Error(codes.InvalidArgument, "id must be greater than 0")
	}

	return nil
}
