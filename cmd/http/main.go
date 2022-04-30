package main

import (
	"log"
	"net/http"
	"os"

	"shorter/internal/service"
	"shorter/internal/storage"
)

func main() {
	db, err := storage.Open()
	if err != nil {
		panic(err)
	}

	s := service.New(db)

	h := http.New()

	http.HandleFunc("/", h)

	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}
