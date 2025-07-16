package app

import (
	"github.com/Braendie/todo-list-storage/internal/config"
	"github.com/Braendie/todo-list-storage/internal/lib/logger"
)

func Start() {

	cfg := config.MustLoad()

	logger := logger.MustLoad(cfg)

	// TODO: реализовать storage

	// TODO: реализовать grpc
}
