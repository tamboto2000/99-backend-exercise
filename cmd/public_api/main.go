package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/tamboto2000/99-backend-exercise/internal/apps/public/app"
	"github.com/tamboto2000/99-backend-exercise/internal/apps/public/config"
	"github.com/tamboto2000/99-backend-exercise/pkg/logger"
)

func main() {
	logger.Info("starting service...")
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal(fmt.Sprintf("error on load config: %v", err))
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())

	a, err := app.RunApp(ctx, cfg)
	if err != nil {
		logger.Fatal(fmt.Sprintf("error on starting app: %v", err))
	}

	logger.Info("service started")

	waitTerminate(a, cancel)
}

func waitTerminate(a *app.App, cancel context.CancelFunc) {
	sigc := make(chan os.Signal, 1)
	signal.Notify(
		sigc,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	<-sigc

	logger.Info("stopping service...")

	cancel()

	if err := a.Stop(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("service stopped")
}
