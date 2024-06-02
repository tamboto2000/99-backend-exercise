package app

import (
	"context"
	"fmt"

	"github.com/tamboto2000/99-backend-exercise/internal/apps/public/config"
	"github.com/tamboto2000/99-backend-exercise/internal/apps/public/modules"
	"github.com/tamboto2000/99-backend-exercise/internal/apps/public/rest"
	"github.com/tamboto2000/99-backend-exercise/pkg/mux"
)

type App struct {
	router *mux.Router
}

func RunApp(_ context.Context, cfg config.Config) (*App, error) {
	// init external services
	exSvcs := modules.NewExternalServices(cfg)

	// init aggregator services
	svcs := modules.NewServices(exSvcs)

	// register routes
	r := mux.NewRouter()
	rest.RegisterRoutes(r, svcs)

	app := App{
		router: r,
	}

	go r.Run(fmt.Sprintf(":%s", cfg.HTTP.Port))

	return &app, nil
}

func (a *App) Stop() error {
	return a.router.Shutdown(context.Background())
}
