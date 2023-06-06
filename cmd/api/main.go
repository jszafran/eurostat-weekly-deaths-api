package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"weekly_deaths/internal/eurostat"

	"github.com/joho/godotenv"
)

// DefaultPort defines a default port that the server will be started on.
const DefaultPort = 8080

func initializeDataSnapshot() (eurostat.DataSnapshot, error) {
	var (
		snapshot eurostat.DataSnapshot
		err      error
	)

	env, ok := os.LookupEnv("DEPLOY_ENV")
	if !ok {
		env = "prod"
	}

	switch env {
	case "local":
		log.Println("DEPLOY_ENV=local; reading snapshot from disk.")
		path := os.Getenv("LOCAL_SNAPSHOT_PATH")
		if path == "" {
			log.Fatal("DEPLOY_ENV set to local and LOCAL_SNAPSHOT_PATH env variable is empty. Exiting.")
		}
		snapshot, err = eurostat.DataSnapshotFromPath(path)
		if err != nil {
			return snapshot, err
		}
	case "production":
		log.Println("DEPLOY_ENV=production; reading live snapshot from Eurostat.")
		snapshot, err = eurostat.DataSnapshotFromEurostat()
		if err != nil {
			return snapshot, err
		}
	}

	return snapshot, nil
}

func main() {
	var port int

	flag.IntVar(&port, "port", DefaultPort, "port to start server on")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found.")
	}

	snapshot, err := initializeDataSnapshot()
	if err != nil {
		log.Fatal(err)
	}

	db := eurostat.DBFromSnapshot(snapshot)
	if err != nil {
		log.Fatal(err)
	}

	app := application{
		db: db,
	}

	router := app.routes()
	log.Printf("Starting the server on :%d port\n", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	if err != nil {
		log.Fatalf("starting server: %s\n", err.Error())
	}
}
