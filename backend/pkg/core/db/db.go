package db

import (
	"fmt"
	"github.com/brissonwilliam/ihavefriends/backend/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Connect to the database and return a handle, or an error
func Connect(cfg config.DB) (*sqlx.DB, error) {
	connectionString := fmt.Sprintf("%s:%s@(%s:%d)/%s?multiStatements=true&parseTime=true", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DatabaseName)
	if cfg.TLS {
		connectionString = connectionString + "&tls=true"
	}
	db, err := sqlx.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
