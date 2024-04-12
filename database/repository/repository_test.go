package repository_test

import (
	"steam-trading/shared/database/model"
	"testing"
	"time"

	"github.com/mikezzb/steam-trading-shared/database"
	"github.com/mikezzb/steam-trading-shared/database/repository"
)

func TestItemRepository_UpdateItem(t *testing.T) {
	t.Run("UpdateItem", func(t *testing.T) {
		db, err := database.NewDBClient("mongodb://localhost:27017", "steam-trading-unit-test", time.Second*10)
		if err != nil {
			t.Error(err)
		}
		db.Ping()
		// defer db.Disconnect()

		repo := repository.NewItemRepository(db)
		item := model.Item{
			Name: "★ Bayonet | Doppler (Factory New)",
		}

		err = repo.UpdateItem(item)

		if err != nil {
			t.Error(err)
		}

		// Get the item back
		updatedItem, err := repo.FindItemByName(item.Name)
		if err != nil {
			t.Error(err)
		}

		if updatedItem == nil {
			t.Errorf("Item not found: %v", item.Name)
		}

	})
}