// Package migration is a simple database migration helper.
// Currently it will only migrate for User Service database
package migration

import (
	"context"
	"fmt"
	"os"

	"github.com/tamboto2000/99-backend-exercise/internal/apps/migration/config"
	"github.com/tamboto2000/99-backend-exercise/pkg/logger"
	"github.com/tamboto2000/99-backend-exercise/pkg/sqli"
)

func MigrateDatabase(cfg config.Config) error {
	for k, c := range cfg.Services {
		err := migrateDb(cfg.Database, c, k)
		if err != nil {
			return err
		}
	}

	return nil
}

func migrateDb(dbCfg config.Database, svcCfg config.Service, svcName string) error {
	log := logger.GetDefault()
	log.WithAttrs(logger.Any("for-service", svcName))

	log.Info(fmt.Sprintf("migrating database for service %s", svcName))

	defDb := "postgres"
	ctx := context.Background()
	db, err := initDatabase(ctx, dbCfg, defDb)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	defer db.Close()

	// read all contents from migration file
	migrateSql, err := os.ReadFile(svcCfg.MigrationFile)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// create the database if not exists
	createDbQ := `
	SELECT 'CREATE DATABASE "%s"'
	WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '%s')	
	`

	createDbQ = fmt.Sprintf(createDbQ, svcCfg.Database, svcCfg.Database)
	var createDbExecQ string
	row := db.QueryRow(ctx, createDbQ)

	err = row.Scan(&createDbExecQ)
	if err != nil {
		if err != sqli.ErrNoRows {
			log.Error(err.Error())
			return err
		}
	}

	if createDbExecQ != "" {
		_, err := db.Exec(ctx, createDbExecQ)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}

	db.Close()

	// after creating the database, create new connection
	// to run the migration file
	db, err = initDatabase(ctx, dbCfg, svcCfg.Database)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	_, err = db.Exec(ctx, string(migrateSql))
	if err != nil {
		log.Error(err.Error())
		return err
	}

	db.Close()

	return nil
}

func initDatabase(ctx context.Context, dbCfg config.Database, dbName string) (*sqli.DB, error) {
	connUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbCfg.Username, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbName)
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
