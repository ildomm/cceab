package main

import (
	"context"
	"errors"
	"github.com/ildomm/cceab/database"
	"github.com/ildomm/cceab/server"
	"github.com/ildomm/cceab/system"
	"log"
	"net/http"
	"os"
)

var (
	gitSha = "unknown" // Populated with the last Git commit SHA (short) at build time
	semVer = "unknown" // Populated with semantic version at build time
)

func main() {
	// Create an overarching context which we can use to safely cancel
	// all goroutines when we receive a signal to terminate.
	ctx, shutdown := context.WithCancel(context.Background())
	defer shutdown()

	// Define standards
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("Server: ")
	log.Printf("starting http server, Version %s, GIT sha %s", semVer, gitSha)

	// Set the timezone to UTC
	system.SetGlobalTimezoneUTC() //nolint:all

	// Parse the command line options
	dBConnURL, err := system.ParseDBConnURL(os.Args[1:])
	if err != nil {
		log.Fatalf("parsing command line: %s", err)
	}
	httpServerPort, err := system.ParseHTTPPort(os.Args[1:])
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

	// Initialize managers
	//userManager := dao.NewUserDAO(querier)
	//gameResultManager := dao.NewGameResultDAO(querier)

	// Initialize the server
	server := server.NewServer()
	server.WithListenAddress(httpServerPort)
	//server.WithUserManager(userManager)
	//server.WithGameResultManager(gameResultManager)

	log.Println("Starting server on", server.ListenAddress())

	if err := server.Run(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("Could not start server on", server.ListenAddress())
		} else {
			log.Println("Server closed")
		}
	}
}
