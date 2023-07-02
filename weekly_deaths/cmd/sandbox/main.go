package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"weekly_deaths/eurostat"
)

func main() {
	err := godotenv.Load("../../.env")
	//err := godotenv.Load("./.env")

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Start fetching latest snapshot")
	sm, err := eurostat.NewSnapshotManager(os.Getenv("S3_BUCKET"))
	if err != nil {
		log.Fatal(err)
	}

	snapshot, err := sm.LatestSnapshotFromS3()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(len(snapshot.Data))
	log.Println(snapshot.Timestamp)
}
