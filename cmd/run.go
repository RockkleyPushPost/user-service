package main

import (
	"context"
	"github.com/RockkleyPushPost/common/config"
	"github.com/RockkleyPushPost/common/database"
	"github.com/RockkleyPushPost/common/di"
	lg "github.com/RockkleyPushPost/common/logger"
	"github.com/RockkleyPushPost/common/setup"
	"github.com/RockkleyPushPost/user-service/internal/metrics"
	"github.com/RockkleyPushPost/user-service/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Config
	cfg, err := config.LoadYamlConfig(os.Getenv("USER_SERVICE_CONFIG_PATH"))

	if err != nil {
		log.Fatal("failed to load config: ", err)
	}

	// Logger
	logger := lg.InitLogger(cfg.ServiceName)

	// Server
	app := setup.NewFiber(fiber.Config{}, cors.Config{}) // fixme

	// Prometheus
	metrics.Init(app, cfg.ServiceName)

	// Database
	db, err := database.NewDatabase(*cfg.Database)

	if err != nil {

		logger.Fatal(err)
	}

	//db.AutoMigrate(&entity.User{})
	//db.AutoMigrate(&entity.Friendship{})
	//db.AutoMigrate(&entity.FriendshipRequest{})

	// DI
	DI := di.NewDI(app, cfg.JwtSecret)

	err = service.Setup(DI, app, db, cfg)

	if err != nil {

		logger.Fatal(err)
	}

	srv, err := service.NewService(
		service.WithConfig(cfg),
		service.WithDI(DI),
		service.WithLogger(logger),
		service.WithServer(app),
	)

	if err != nil {

		logger.Fatal(err)
	}

	go handleShutdown(ctx, cancel, srv, logger)

	logger.Fatal(srv.Run(ctx))

}

func handleShutdown(ctx context.Context, cancel context.CancelFunc, srv service.Service, logger *log.Logger) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		logger.Printf("received signal: %v", sig)
		cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logger.Printf("shutdown error: %v", err)
		}
	case <-ctx.Done():
		logger.Println("context cancelled")
	}
}
