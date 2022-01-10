package storage

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type DB struct {
	*sqlx.DB
}

type Data struct {
	ID    int
	URL   string
	Short string
}

func Open() (*DB, error) {
	name := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOSTNAME")
	dbName := os.Getenv("DB_NAME")
	mode := os.Getenv("DB_MODE")

	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		name, password, host, dbName, mode)

	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	log.Println("connected")

	return &DB{DB: db}, nil
}

func (db *DB) Insert(URL, short string) error {
	_, err := db.Exec("INSERT INTO links (url, short) VALUES ($1, $2)", URL, short)
	return err
}

func (db *DB) ByShort(short string) (Data, error) {
	var data Data
	row := db.QueryRow("SELECT * FROM links WHERE short=($1)", short)
	err := row.Scan(&data.ID, &data.URL, &data.Short)
	return data, err
}

func (db *DB) ByURL(url string) (Data, error) {
	var data Data
	row := db.QueryRow("SELECT * FROM links WHERE url=($1)", url)
	err := row.Scan(&data.ID, &data.URL, &data.Short)
	return data, err
}
