package database_test

import (
	"log"
	"testing"
	"time"

	"github.com/mikezzb/steam-trading-shared/database"
)

func TestDBClient(t *testing.T) {
	localMongoURL := "mongodb://localhost:27017"
	dbName := "steam-trading"
	t.Run("Connect", func(t *testing.T) {
		dbClient, err := database.NewDBClient(localMongoURL, dbName, 10*time.Second)
		if err != nil {
			t.Fatalf("Failed to connect to db: %v", err)
		}
		log.Printf("Connected to db at %s\n", localMongoURL)
		defer dbClient.Disconnect()
	})
}

func TestDBRepositories(t *testing.T) {
	localMongoURL := "mongodb://localhost:27017"
	dbName := "steam-trading"
	dbClient, err := database.NewDBClient(localMongoURL, dbName, 10*time.Second)
	t.Run("NewRepositories", func(t *testing.T) {
		if err != nil {
			t.Fatalf("Failed to connect to db: %v", err)
		}
		repos := database.NewRepositories(dbClient)
		log.Printf("Repositories: %v\n", repos)
		defer dbClient.Disconnect()
	})
}
