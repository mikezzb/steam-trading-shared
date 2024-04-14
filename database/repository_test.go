package database_test

import (
	"log"
	"testing"
	"time"

	"github.com/mikezzb/steam-trading-shared/database/model"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/mikezzb/steam-trading-shared/database"
)

func TestItemRepository_UpdateItem(t *testing.T) {
	t.Run("UpdateItem", func(t *testing.T) {
		db, err := database.NewDBClient("mongodb://localhost:27017", "steam-trading-unit-test", time.Second*10)
		if err != nil {
			t.Error(err)
		}
		defer db.Disconnect()
		repos := database.NewRepositories(db)

		repo := repos.GetItemRepository()
		item := &model.Item{
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
			return
		}

		// delete
		err = repo.DeleteItemByName(updatedItem)

		if err != nil {
			t.Error(err)
		}

		_, err = repo.FindItemByName(item.Name)
		if err == nil {
			t.Errorf("Item not deleted: %v", item.Name)
		}

	})
}

func TestListingRepo_Insert(t *testing.T) {
	t.Run("Insert", func(t *testing.T) {
		db, err := database.NewDBClient("mongodb://localhost:27017", "steam-trading-unit-test", time.Second*10)
		if err != nil {
			t.Error(err)
		}
		defer db.Disconnect()
		repos := database.NewRepositories(db)

		repo := repos.GetListingRepository()
		// listingFile := "mocks/listings.json"
		// b, err := os.ReadFile(listingFile)
		if err != nil {
			t.Error(err)
		}

		listings := []model.Listing{
			{
				Name: "★ Bayonet | Doppler (Factory New)",
			},
			{
				Name: "★ Bayonet | Doppler (Minimal Wear)",
			},
		}
		// json.Unmarshal(b, &listings)

		if err != nil {
			t.Error(err)
		}

		err = repo.InsertListings(listings)

		if err != nil {
			t.Error(err)
		}

		// Get the item back
		updatedListing, err := repo.FindListingByItemName(listings[0].Name)
		if err != nil {
			t.Error(err)
		}
		log.Printf("Updated listing: %v", updatedListing)

		if updatedListing == nil {
			t.Errorf("Listing not found: %v", listings[0].Name)
			return
		}

		// delete
		err = repo.DeleteListingByItemName(updatedListing.Name)
		if err != nil {
			t.Error(err)
		}
	})

}

func TestTransactionRepo_Insert(t *testing.T) {
	t.Run("Insert", func(t *testing.T) {
		db, err := database.NewDBClient("mongodb://localhost:27017", "steam-trading-unit-test", time.Second*10)
		if err != nil {
			t.Error(err)
		}
		defer db.Disconnect()
		repos := database.NewRepositories(db)

		repo := repos.GetTransactionRepository()
		transactions := []model.Transaction{
			{
				Name: "★ Bayonet | Doppler (Factory New)",
			},
			{
				Name: "★ Bayonet | Doppler (Minimal Wear)",
			},
		}

		err = repo.InsertTransactions(transactions)

		if err != nil {
			t.Error(err)
		}

		// Get the item back
		updatedTransaction, err := repo.FindTransactionByItemName(transactions[0].Name)
		if err != nil {
			t.Error(err)
		}

		if updatedTransaction == nil {
			t.Errorf("Transaction not found: %v", transactions[0].Name)
			return
		}

		// delete
		err = repo.DeleteTransactionByItemName(updatedTransaction.Name)
		if err != nil {
			t.Error(err)
		}
	})
}

func TestMongoID(t *testing.T) {
	t.Run("MongoID", func(t *testing.T) {

		// Generate a new ObjectID
		objectID := primitive.NewObjectID()

		// Extract timestamp information from the ObjectID
		timestamp := objectID.Timestamp()
		log.Printf("Timestamp: %v\n", timestamp)

	})
}
