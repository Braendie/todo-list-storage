package grpc

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/Braendie/todo-list-storage/internal/grpc/storage"
	storageService "github.com/Braendie/todo-list-storage/internal/services/storage"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
	address    string
}

func New(
	log *slog.Logger,
	storageService storageService.StorageService,
	port int,
	address string,
) *App {
	gRPCServer := grpc.NewServer()
	storage.Register(gRPCServer, &storageService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
		address:    address,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "app.grpc.Run"

	log := a.log.With(slog.String("op", op))

	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.address, a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grpc server is running", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}
