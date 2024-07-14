package main

import (
	"context"
	"github.com/ildomm/cceab/dao"
	"github.com/ildomm/cceab/database"
	"github.com/ildomm/cceab/system"
	"log"
	"os"
	"time"
)

var (
	gitSha = "unknown" // Populated with the last Git commit SHA (short) at build time
	semVer = "unknown" // Populated with semantic version at build time

	pauseDuration      = 1 * time.Minute
	totalGamesToCancel = 10
)

func main() {
	// Create an overarching context which we can use to safely cancel
	// all goroutines when we receive a signal to terminate.
	ctx, shutdown := context.WithCancel(context.Background())
	defer shutdown()

	// Define standards
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("Validator: ")
	log.Printf("starting validator, Version %s, GIT sha %s", semVer, gitSha)

	// Set the timezone to UTC
	system.SetGlobalTimezoneUTC() //nolint:all

	// Parse the command line options
	dBConnURL, err := system.ParseDBConnURL(os.Args[1:])
	if err != nil {
		log.Fatalf("parsing command line: %s", err)
	}

	// Set up the database connection and run migrations
	log.Printf("connecting to database")
	querier, err := database.NewPostgresQuerier(
		ctx,
		dBConnURL,
	)
	if err != nil {
		log.Fatalf("error connecting to the database: %s", err)
	}
	defer querier.Close()

	// Initialize manager
	gameResultManager := dao.NewGameResultDAO(querier)

	// Run the pipeline
	go run(ctx, gameResultManager)

	log.Printf("caught signal, terminating. %v", system.WaitForSignal().String())
}

// Run starts the validation pipeline
func run(ctx context.Context, gameResultManager dao.GameResultDAO) {
	log.Printf("starting validating game results")

	for {
		select {
		case <-ctx.Done():
			log.Printf("stopping pipeline")
			return
		default:
			err := gameResultManager.ValidateGameResults(ctx, totalGamesToCancel)
			if err != nil {
				log.Printf("error validating game results: %s", err)
			}

			system.SleepWithContext(ctx, pauseDuration)
		}
	}
}
