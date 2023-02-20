package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	// SQLite driver
	_ "github.com/mattn/go-sqlite3"

	"github.com/3d0c/storage/pkg/config"
)

var (
	instance *sql.DB
	once     sync.Once
)

// Instance is a DB connection singleton
func Instance(cfg config.Database) *sql.DB {
	once.Do(func() {
		var err error

		if instance, err = connectDatabase(cfg); err != nil {
			panic(err)
		}
	})

	return instance
}

func connectDatabase(cfg config.Database) (*sql.DB, error) {
	var (
		conn *sql.DB
		dsn  = cfg.DSN
		err  error
	)

	if dsn == "" {
		return nil, fmt.Errorf("error opening database - DataSource can't be empty")
	}

	if err = os.MkdirAll(filepath.Dir(dsn), os.ModePerm); err != nil {
		return nil, err
	}

	if conn, err = sql.Open(cfg.Dialect, dsn); err != nil {
		return nil, err
	}

	if err = migrate(conn, cfg); err != nil {
		return nil, fmt.Errorf("error migrating database - %s", err)
	}

	return conn, nil
}
