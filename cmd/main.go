package main

import (
	"net/http"
	"os"

	"github.com/hararudoka/shrt/handler"
	"github.com/hararudoka/shrt/service"
	"github.com/hararudoka/shrt/storage"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	db, err := storage.Open()
	if err != nil {
		log.Fatal().Err(err)
	}
	log.Info().Msg("database connection established")

	s := service.New(db)

	handler := handler.New(*s)

	if err := http.ListenAndServe(":"+os.Getenv("PORT"), handler); err != nil {
		log.Fatal().Err(err).Msg("error during ListenAndServe")
	}
}
