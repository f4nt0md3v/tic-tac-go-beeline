package postgres

import (
	"database/sql"
	"strconv"

	_ "github.com/lib/pq"

	"github.com/f4nt0md3v/tic-tac-go-beeline/app/pkg/env"
)

const (
	migrateScript = `
		CREATE TABLE IF NOT EXISTS games(
		    id                SERIAL UNIQUE NOT NULL PRIMARY KEY,
		    game_id           VARCHAR (50) NOT NULL,
		    first_user_id     VARCHAR (50),
		    second_user_id    VARCHAR (50),
		    state             VARCHAR (50) NOT NULL,
		    last_move_user_id VARCHAR (50) NOT NULL,
		    created_at        TIMESTAMP NOT NULL DEFAULT NOW(),
		    last_modified_at  TIMESTAMP NOT NULL DEFAULT NOW()
		);
	`
)

type Config struct {
	Host             string
	Port             int
	User             string
	Password         string
	Database         string
	Params           string
	ConnectionString string
	Mode             env.Mode
}

func NewDBSession(cfg Config) (*sql.DB, error) {
	return getDbConn(getDbConnString(cfg))
}

func getDbConnString(cfg Config) string {
	var connStr string
	if cfg.ConnectionString == "" {
		connStr = "postgres://"
		if cfg.Host == "" {
			cfg.Host = "localhost"
		}
		if cfg.Port == 0 {
			cfg.Port = 5432
		}
		if cfg.Params == "" {
			if cfg.Mode == env.Development {
				cfg.Params = "sslmode=disable"
			}
		}
		connStr +=
			cfg.User + ":" +
				cfg.Password + "@" +
				cfg.Host + ":" +
				strconv.Itoa(cfg.Port) + "/" +
				cfg.Database + "?" +
				cfg.Params
	} else {
		connStr = cfg.ConnectionString
	}
	return connStr
}

func getDbConn(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Migrate(db *sql.DB) error {
	_, err := db.Exec(migrateScript)
	return err
}
