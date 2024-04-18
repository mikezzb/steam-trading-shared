package subscription_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	"github.com/mikezzb/steam-trading-shared/database"
	"github.com/mikezzb/steam-trading-shared/database/model"
	"github.com/mikezzb/steam-trading-shared/database/repository"
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

		repos := repository.NewRepoFactory(dbClient, nil)

		// create item
		itemRepo := repos.GetItemRepository()
		item := &model.Item{
			Name:              "★ Bayonet | Marble Fade (Factory New)",
			LowestMarketPrice: "1000",
		}
		itemRepo.UpdateItem(item)

		// create subs
		subRepo := repos.GetSubscriptionRepository()
		sub := &model.Subscription{
			Name:       "★ Bayonet | Marble Fade (Factory New)",
			MaxPremium: "1.0",
			Rarity:     "FFI",
			NotiType:   "telegram",
			NotiId:     secrets["telegramTestChatId"],
		}
		subRepo.UpsertSubscription(sub)

		emitter := subscription.NewNotificationEmitter(
			&subscription.NotifierConfig{
				TelegramToken: secrets["telegramToken"],
			},
		)

		emitter.Init(repos)

		emitter.EmitListing(&model.Listing{
			Name:   "★ Bayonet | Marble Fade (Factory New)",
			Rarity: "FFI",
			Price:  "1001",
		})

		time.Sleep(1 * time.Second)

		subRepo.DeleteSubscriptionByName(
			sub.Name,
		)

		itemRepo.DeleteAll()

	})
}
