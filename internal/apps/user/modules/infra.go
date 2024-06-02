package modules

import (
	"context"

	"github.com/tamboto2000/99-backend-exercise/internal/apps/user/config"
	"github.com/tamboto2000/99-backend-exercise/internal/apps/user/infra"
	"github.com/tamboto2000/99-backend-exercise/pkg/sqli"
)

type Infra struct {
	DB *sqli.DB
}

func NewInfra(ctx context.Context, cfg config.Config) (inf *Infra, err error) {
	// snowid
	if err := infra.InitSnowID(); err != nil {
		return nil, err
	}

	// database
	db, err := infra.InitDatabase(ctx, cfg.Database)
	if err != nil {
		return
	}

	inf = &Infra{
		DB: db,
	}

	return inf, nil
}
