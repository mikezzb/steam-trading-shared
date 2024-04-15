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
	db, err := database.NewDBClient("mongodb://localhost:27017", "steam-trading-unit-test", time.Second*10)
	if err != nil {
		t.Error(err)
	}
	defer db.Disconnect()
	repos := database.NewRepositories(db)

	repo := repos.GetTransactionRepository()

	t.Run("Insert", func(t *testing.T) {
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

	t.Run("Upsert", func(t *testing.T) {
		transactions := []model.Transaction{
			{
				Name:       "★ Bayonet | Doppler (Factory New)",
				AssetId:    "123",
				PreviewUrl: "Old Preview URL",
			},
			{
				Name:    "★ Bayonet | Doppler (Minimal Wear)",
				AssetId: "456",
			},
		}

		err = repo.InsertTransactions(transactions)

		if err != nil {
			t.Error(err)
		}

		// Upsert the transaction
		newTransactions := []model.Transaction{
			{
				Name:       "★ Bayonet | Doppler (Factory New)",
				AssetId:    "123",
				PreviewUrl: "New Preview URL",
			},
			{
				Name:    "★ Bayonet | Doppler (Minimal Wear)",
				AssetId: "101112",
			},
		}

		err = repo.UpsertTransactionsByAssetID(newTransactions)

		if err != nil {
			t.Error(err)
		}

		// Get the item back
		updatedTransaction, err := repo.FindTransactionByAssetId(transactions[0].AssetId)
		if err != nil {
			t.Error(err)
		}

		// expect the preview url to be updated
		if updatedTransaction.PreviewUrl != newTransactions[0].PreviewUrl {
			t.Errorf("Preview URL not updated: %v", updatedTransaction.PreviewUrl)
		}

		// get ALL transactions
		allTransactions, err := repo.FindAllTransactions()
		if err != nil {
			t.Error(err)
		}

		log.Printf("%v transactions: %v", len(allTransactions), allTransactions)

		// delete all transactions
		err = repo.DeleteAll()
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

func TestDeleteOldListing(t *testing.T) {
	db, err := database.NewDBClient("mongodb://localhost:27017", "steam-trading-unit-test", time.Second*10)
	if err != nil {
		t.Error(err)
	}
	defer db.Disconnect()
	repos := database.NewRepositories(db)

	repo := repos.GetListingRepository()

	t.Run("DeleteOldListing", func(t *testing.T) {
		listings := []model.Listing{
			{
				Name:    "★ Bayonet | Doppler (Factory New)",
				AssetId: "123",
				Price:   "100",
			},
			{
				Name:    "★ Bayonet | Doppler (Factory New)",
				AssetId: "123",
				Price:   "105", // UPDATED
			},
		}

		err = repo.InsertListings(listings)

		if err != nil {
			t.Error(err)
		}

		err = repo.DeleteOldListingsByAssetID(listings[0].AssetId)
		if err != nil {
			t.Error(err)
		}

		// Get the item back
		updatedListing, err := repo.FindItemByAssetId(listings[0].AssetId)
		if err != nil {
			t.Error(err)
		}

		// ensure the list have updated price
		if updatedListing.Price != listings[1].Price {
			t.Errorf("Price not updated: %v", updatedListing.Price)
		}
	})
}

func TestMergeByAssetIds(t *testing.T) {
	db, err := database.NewDBClient("mongodb://localhost:27017", "steam-trading-unit-test", time.Second*10)
	if err != nil {
		t.Error(err)
	}
	defer db.Disconnect()
	repos := database.NewRepositories(db)

	listingRepo := repos.GetListingRepository()
	transactionRepo := repos.GetTransactionRepository()

	t.Run("Merge Listings", func(t *testing.T) {
		assetIds, err := listingRepo.GetAllUniqueAssetIDs()
		if err != nil {
			t.Error(err)
		}

		for _, assetId := range assetIds {
			log.Printf("Merging assetId: %v\n", assetId)
			listingRepo.DeleteOldListingsByAssetID(assetId)
		}
	})

	t.Run("Merge Transactions", func(t *testing.T) {
		assetIds, err := transactionRepo.GetAllUniqueAssetIDs()
		if err != nil {
			t.Error(err)
		}

		for _, assetId := range assetIds {
			log.Printf("Merging assetId: %v\n", assetId)
			transactionRepo.DeleteOldTransactionsByAssetID(assetId)
		}
	})
}
