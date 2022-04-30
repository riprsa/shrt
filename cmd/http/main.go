package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"shorter/internal/service"
	"shorter/internal/storage"
)

func main() {
	mode := ""
	fmt.Println("choose mode:\n1. without pem\n2. with pem")
	fmt.Scan(&mode)

	db, err := storage.Open()
	if err != nil {
		panic(err)
	}

	s := service.New(db)

	http.HandleFunc("/", h)

	if mode == "1" {
		if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
			log.Fatal(err)
		}
	} else {
		// pem := os.Getenv("SERVER_PEM")
		// key := os.Getenv("SERVER_KEY")

		// if pem == "" || key == "" {
		// 	panic("pem or key is empty")
		// }

		// if err := e.StartTLS(":"+os.Getenv("PORT"), []byte(pem), []byte(key)); err != http.ErrServerClosed {
		// 	panic(err)
		// }
	}
}
