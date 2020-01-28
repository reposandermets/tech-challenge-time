package models

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	host     = "0.0.0.0"
	port     = 5432
	user     = "api_user"
	password = "123456"
	dbname   = "reports"
)

// Connect ...
func Connect() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
