package main

import (
	"os"
	"rest-api/gateways"
	"rest-api/usecases"
	"time"

	"github.com/rs/zerolog"
)

func main() {
	db := usecases.NewDatabase()

	logger := zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
	).Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	logger.Info().Msg("Server started")

	api := gateways.NewBooksAPI(db, &logger)
	api.Handle()

	err := api.Serve()
	if err != nil {
		panic(err)
	}
}
