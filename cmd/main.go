package main

import (
	"context"
	"os"
	"os/signal"
	"store-service/config"
	"store-service/utils"
	"syscall"

	"github.com/Slicky-Team/slickfame/sqlengine"
	"github.com/rs/zerolog/log"

	"golang.org/x/sync/errgroup"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	// Connect to PostgreSQL
	dbSource := cfg.DatabaseURL
	dbEngine, err := sqlengine.NewPostgresDB(sqlengine.DBConnString(dbSource))
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to connect to database: %v\n", err)
	}
	defer dbEngine.Close()

	// Run migrations
	err = utils.RunMigrations(ctx, dbEngine)
	if err != nil {
		panic(err)
	}

	waitGroup, ctx := errgroup.WithContext(ctx)

	err = waitGroup.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("Error from wait group")
	}

}
