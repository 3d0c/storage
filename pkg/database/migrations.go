package database

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"

	"github.com/3d0c/storage/pkg/config"
)

//go:embed migrations/*.sql
var migrations embed.FS

// Migrate migrates database scheme
func migrate(conn *sql.DB) error {
	var (
		err error
	)

	goose.SetBaseFS(migrations)

	if err = goose.SetDialect(config.Proxy().Database.Dialect); err != nil {
		return err
	}

	if err = goose.Up(conn, "migrations"); err != nil {
		return err
	}

	return nil
}
