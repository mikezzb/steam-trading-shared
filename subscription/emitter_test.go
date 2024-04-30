package subscription_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	shared "github.com/mikezzb/steam-trading-shared"
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
		items := []*model.Item{
			{
				Name: "★ Bayonet | Marble Fade (Factory New)",
				BuffPrice: &model.MarketPrice{
					Price:     shared.GetDecimal128("1000"),
					UpdatedAt: time.Now(),
				},
			},
			{
				Name: "★ Flip Knife | Marble Fade (Factory New)",
				BuffPrice: &model.MarketPrice{
					Price:     shared.GetDecimal128("100"),
					UpdatedAt: time.Now(),
				},
			},
		}

		for _, item := range items {
			itemRepo.UpsertItem(item)
		}

		// create subs
		subRepo := repos.GetSubscriptionRepository()
		subs := []*model.Subscription{
			{
				Name:       "★ Bayonet | Marble Fade (Factory New)",
				MaxPremium: "1.0",
				Rarity:     "FFI",
				NotiType:   "telegram",
				NotiId:     secrets["telegramTestChatId"],
			},
			{
				Name:       "★ Flip Knife | Marble Fade (Factory New)",
				MaxPremium: "1.0",
				Rarity:     "FFI",
				NotiType:   "telegram",
				NotiId:     secrets["telegramTestChatId"],
			},
		}

		for _, sub := range subs {
			subRepo.InsertSubscription(sub)
		}

		time.Sleep(500 * time.Millisecond)

		emitter := subscription.NewNotificationEmitter(
			&subscription.NotifierConfig{
				TelegramToken: secrets["telegramToken"],
			},
		)

		emitter.Init(repos)

		emitter.EmitListing(&model.Listing{
			Name:   "★ Bayonet | Marble Fade (Factory New)",
			Rarity: "FFI",
			Price:  shared.GetDecimal128("1001"),
		})

		time.Sleep(1 * time.Second)

		subRepo.DeleteAll()

		itemRepo.DeleteAll()

	})
}
