package main

import (
	"flag"
	"os"
	"time"

	"github.com/kgantsov/filo/internal/db"
	"github.com/kgantsov/filo/internal/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	logLevel string
	port     int
)

func main() {
	// Initialize database
	filoDB := db.InitDB("filo.db")

	flag.StringVar(&logLevel, "log_level", "debug", "Log level")
	flag.IntVar(&port, "port", 8000, "Server port")
	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339Nano})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixNano

	logLevel, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(logLevel)
	}

	srv := server.NewServer(filoDB, port)

	log.Info().Msgf("Server starting on http://localhost:%d", port)
	log.Info().Msgf("API docs available at http://localhost:%d/docs", port)

	if err := srv.Start(); err != nil {
		log.Fatal().Msgf("Server failed to start: %v", err)
	}
}
