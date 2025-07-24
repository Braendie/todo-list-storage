package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Braendie/todo-list-storage/internal/app/grpc"
	"github.com/Braendie/todo-list-storage/internal/config"
	"github.com/Braendie/todo-list-storage/internal/lib/logger"
	storageService "github.com/Braendie/todo-list-storage/internal/services/storage"
	"github.com/Braendie/todo-list-storage/internal/storage/pgsql"
)

func Start() {

	cfg := config.MustLoad()

	logger := logger.MustLoad(cfg)

	storage := pgsql.New(cfg.StoragePath)

	storageService := storageService.New(logger, storage)

	server := grpc.New(logger, *storageService, cfg.Server.Port, cfg.Server.Address)

	go server.MustRun()

	logger.Info("Server started")

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGTERM, os.Interrupt)

	<-stop

	logger.Info("Server stopped")

	server.Stop()
}
