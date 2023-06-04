package storage

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	*pgxpool.Pool
}

func Open() (*DB, error) {
	name := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOSTNAME")
	dbName := os.Getenv("DB_NAME")
	mode := os.Getenv("DB_MODE")
	port := os.Getenv("DB_PORT")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", name, password, host, port, dbName, mode)

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	// defer conn.Close(context.Background())
	return &DB{pool}, nil
}

func (db *DB) ByShort(short string) (string, error) {
	var url string
	row := db.QueryRow(context.Background(), "SELECT url FROM links WHERE short=$1", short)
	err := row.Scan(&url)
	return url, err
}

func (db *DB) ByURL(url string) (string, error) {
	var short string
	row := db.QueryRow(context.Background(), "SELECT short FROM links WHERE url=$1", url)
	err := row.Scan(&short)
	return short, err
}

func (db *DB) Insert(URL, short string) error {
	_, err := db.Exec(context.Background(), "INSERT INTO links (url, short) VALUES ($1, $2)", URL, short)
	return err
}
