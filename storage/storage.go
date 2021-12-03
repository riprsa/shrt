package storage

import (
	"fmt"
	"log"
	"github.com/jmoiron/sqlx"
	"net/url"
	"shorter/config"

	_ "github.com/lib/pq"
)

type DB struct {
	*sqlx.DB
}

type Data struct {
	ID int
	URL string
	Short string
}

func Open() (*DB, error) {
	var conf config.Config
	err := conf.LoadData()
	if err != nil {
		return nil, err
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s%s",
		conf.Username, conf.Password, conf.Host, conf.DBName, conf.Args)
	fmt.Println(connStr)
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	log.Println("connected")

	return &DB{DB:db}, nil
}

func (db *DB) Insert(URL, short string) error {
	_, err := url.Parse(URL)
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO links (url, short) VALUES ($1, $2)", URL, short)
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
	err := row.Scan(&data.ID,&data.URL, &data.Short)
	return data, err
}
