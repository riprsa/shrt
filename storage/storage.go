package storage

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

type DB struct {
	*sqlx.DB
}

func Open() (*DB, error) {
	name := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOSTNAME")
	dbName := os.Getenv("DB_NAME")
	mode := os.Getenv("DB_MODE")

	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		name, password, host, dbName, mode)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &DB{DB: db}, nil
}

func (db *DB) Insert(URL, short string) error {
	_, err := db.Exec("INSERT INTO links (url, short) VALUES ($1, $2)", URL, short)
	return err
}

func (db *DB) ByShort(short string) (string, error) {
	row := db.QueryRow("SELECT url FROM links WHERE short=($1)", short)
	var url string
	err := row.Scan(&url)
	return url, err
}

func (db *DB) ByURL(url string) (string, error) {
	row := db.QueryRow("SELECT short FROM links WHERE url=($1)", url)
	var short string
	err := row.Scan(&short)
	return short, err
}
