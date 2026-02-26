package app

import (
	"context"
	"embed"

	"github.com/doug-martin/goqu/v9"
	pgdb "github.com/itmo-lite-chat/go-utils/postgres_db"
	"github.com/pkg/errors"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func (a *app) initPostgresDB(ctx context.Context) (*goqu.Database, error) {
	dbConfig := a.config.DBConfig

	connString := pgdb.GetConnectionString(
		dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name,
	)

	dbconn, err := pgdb.Connect(ctx, connString, dbConfig.ConnectionTTL, dbConfig.MaxOpenConnections, dbConfig.MaxIdleConnections)
	if err != nil {
		return nil, errors.Wrap(err, "can't connect to pgdb")
	}
	a.closer.AddCloseFunc("postgres", dbconn.Close)

	if a.config.DBConfig.ApplyMigration {
		if err := pgdb.ApplyMigrations(ctx, dbconn, embedMigrations); err != nil {
			return nil, errors.Wrap(err, "can't apply migrations")
		}
	}

	return goqu.New("postgres", dbconn), nil
}
