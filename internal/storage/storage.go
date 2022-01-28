package storage

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
	"shorter/internal/model"

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

	log.Println(connStr)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}

	log.Println("connected")

	return &DB{DB: db}, nil
}

func (db *DB) Insert(URL, short string) error {
	_, err := db.Exec("INSERT INTO links (url, short) VALUES ($1, $2)", URL, short)
	return err
}

func (db *DB) ByShort(short string) (model.Data, error) {
	var data model.Data
	row := db.QueryRow("SELECT * FROM links WHERE short=($1)", short)
	err := row.Scan(&data.ID, &data.URL, &data.Short)
	return data, err
}

func (db *DB) ByURL(url string) (model.Data, error) {
	var data model.Data
	row := db.QueryRow("SELECT * FROM links WHERE url=($1)", url)
	err := row.Scan(&data.ID, &data.URL, &data.Short)
	return data, err
}
