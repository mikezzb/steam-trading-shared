package subscription_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	"github.com/mikezzb/steam-trading-shared/database"
	"github.com/mikezzb/steam-trading-shared/database/model"
	"github.com/mikezzb/steam-trading-shared/subscription"
)

func TestEmitter(t *testing.T) {
	// load secrets from json
	secretsPath := "./secrets.json"
	secretsFile, err := os.ReadFile(secretsPath)
	if err != nil {
		log.Fatalf("TestTelegram: %v", err)
	}
	var secrets map[string]string
	// unmarshal json
	json.Unmarshal(secretsFile, &secrets)

	t.Run("TestEmitter", func(t *testing.T) {
		dbClient, err := database.NewDBClient(
			"mongodb://localhost:27017",
			"steam-trading-unit-test",
			10*time.Second,
		)
		defer dbClient.Disconnect()
		if err != nil {
			t.Fatalf("TestEmitter: %v", err)
		}

		repos := database.NewRepositories(dbClient)

		emitter := subscription.NewNotificationEmitter(
			repos.GetSubscriptionRepository(),
			repos.GetItemRepository(),
			&subscription.NotifierConfig{
				TelegramToken: secrets["telegramToken"],
			},
		)

		emitter.EmitListing(&model.Listing{
			Name:   "★ Bayonet | Marble Fade (Factory New)",
			Rarity: "FFI",
			Price:  "-1",
		})

		time.Sleep(1 * time.Second)

	})
}
