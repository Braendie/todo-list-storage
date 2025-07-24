package app

import (
	"github.com/Braendie/todo-list-storage/internal/config"
	"github.com/Braendie/todo-list-storage/internal/lib/logger"
	"github.com/Braendie/todo-list-storage/internal/storage/pgsql"
)

func Start() {

	cfg := config.MustLoad()

	logger := logger.MustLoad(cfg)

	storage := pgsql.New(cfg.StoragePath)

	// TODO: реализовать grpc
}
