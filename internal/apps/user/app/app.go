package app

import (
	"context"
	"fmt"

	"github.com/tamboto2000/99-backend-exercise/internal/apps/user/config"
	"github.com/tamboto2000/99-backend-exercise/internal/apps/user/modules"
	"github.com/tamboto2000/99-backend-exercise/internal/apps/user/rest"
	"github.com/tamboto2000/99-backend-exercise/pkg/mux"
)

type App struct {
	inf    *modules.Infra
	router *mux.Router
}

func RunApp(ctx context.Context, cfg config.Config) (*App, error) {
	// init infra
	inf, err := modules.NewInfra(ctx, cfg)
	if err != nil {
		return nil, err
	}

	// init services
	svcs := modules.NewServices(inf)

	// register REST routes
	r := mux.NewRouter()
	rest.RegisterRoutes(r, svcs)

	// run router
	go r.Run(fmt.Sprintf(":%s", cfg.HTTP.Port))

	app := App{
		inf:    inf,
		router: r,
	}

	return &app, nil
}

func (a *App) Stop() error {
	if err := a.inf.DB.Close(); err != nil {
		return err
	}

	return a.router.Shutdown(context.Background())
}
