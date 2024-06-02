package infra

import (
	"context"
	"fmt"

	"github.com/tamboto2000/99-backend-exercise/internal/apps/user/config"
	"github.com/tamboto2000/99-backend-exercise/pkg/sqli"
)

func InitDatabase(ctx context.Context, dbCfg config.Database) (*sqli.DB, error) {
	connUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbCfg.Username, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.Database)
	conn, err := sqli.NewPostgreConn(ctx, connUrl)
	if err != nil {
		return nil, err
	}

	db := sqli.NewDB(conn)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
