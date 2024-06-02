package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/tamboto2000/99-backend-exercise/internal/apps/migration"
	"github.com/tamboto2000/99-backend-exercise/internal/apps/migration/config"
	"github.com/tamboto2000/99-backend-exercise/pkg/logger"
)

func main() {
	logger.Info("begin database migration...")

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal(err.Error())
	}

	if err := migration.MigrateDatabase(cfg); err != nil {
		logger.Fatal(err.Error())
	}

	logger.Info("database migration finished")

	// only exit when receive termination signals
	// so that this program will not boot-loop
	waitTerminate()
}

func waitTerminate() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(
		sigc,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	<-sigc
}
