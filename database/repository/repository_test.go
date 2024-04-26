package repository_test

import (
	"log"
	"testing"
	"time"

	"github.com/mikezzb/steam-trading-shared/database/model"
	"github.com/mikezzb/steam-trading-shared/database/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/mikezzb/steam-trading-shared/database"
)

func RepoInit() (*database.DBClient, repository.RepoFactory, error) {
	db, err := database.NewDBClient("mongodb://localhost:27017", "steam-trading-unit-test", time.Second*10)
	if err != nil {
		log.Fatal(err)
	}
	return db, repository.NewRepoFactory(db, nil), nil
}

func TestItemRepository_UpdateItem(t *testing.T) {
	t.Run("UpdateItem", func(t *testing.T) {
		db, repos, err := RepoInit()
		defer db.Disconnect()
		if err != nil {
			t.Error(err)
		}

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
	db, err := database.NewDBClient("mongodb://localhost:27017", "steam-trading-unit-test", time.Second*10)
	if err != nil {
		t.Error(err)
	}
	defer db.Disconnect()
	repos := repository.NewRepoFactory(db, nil)
	repo := repos.GetListingRepository()

	t.Run("Insert", func(t *testing.T) {

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
		updatedListing, err := repo.GetListingByItemName(listings[0].Name)
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

func TestListingRepo_Upsert(t *testing.T) {
	db, repos, err := RepoInit()
	defer db.Disconnect()
	if err != nil {
		t.Error(err)
	}
	repo := repos.GetListingRepository()

	t.Run("Upsert", func(t *testing.T) {
		listings := []model.Listing{
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

		err = repo.InsertListings(listings)

		if err != nil {
			t.Error(err)
		}

		// Upsert the transaction
		newListings := []model.Listing{
			// Shall update
			{
				Name:       "★ Bayonet | Doppler (Factory New)",
				AssetId:    "123",
				PreviewUrl: "New Preview URL",
			},
			// Shall insert
			{
				Name:    "★ Bayonet | Doppler (Minimal Wear)",
				AssetId: "101112",
			},
			// Shall NOT update
			{
				Name:    "★ Bayonet | Doppler (Minimal Wear)",
				AssetId: "456",
			},
		}

		updatedListings, err := repo.UpsertListingsByAssetID(newListings)

		if err != nil {
			t.Error(err)
		}

		log.Printf("Updated listings: %v", updatedListings)

		// Get the item back
		updatedListing, err := repo.FindOneListing(bson.M{
			"assetId": listings[0].AssetId,
		})
		if err != nil {
			t.Error(err)
		}

		log.Printf("Updated listing: %v", updatedListing)

		// expect the preview url to be updated
		if updatedListing.PreviewUrl != updatedListings[0].PreviewUrl {
			t.Errorf("Preview URL not updated: %v", updatedListing.PreviewUrl)
		}

		// delete all transactions
		err = repo.DeleteAll()
		if err != nil {
			t.Error(err)
		}
	})

}

func TestTransactionRepo_Insert(t *testing.T) {
	db, repos, err := RepoInit()
	defer db.Disconnect()
	if err != nil {
		t.Error(err)
	}

	repo := repos.GetTransactionRepository()

	t.Run("Insert", func(t *testing.T) {
		transactions := []model.Transaction{
			{
				Name:      "★ Bayonet | Doppler (Factory New)",
				CreatedAt: time.Now(),
			},
			{
				Name:      "★ Bayonet | Doppler (Minimal Wear)",
				CreatedAt: time.Now(),
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
				PreviewUrl: "Old Preview URL",
				CreatedAt:  time.Now(),
				Metadata: model.TransactionMetadata{
					AssetId: "123",
					Market:  "buff",
				},
			},
			{
				Name: "★ Bayonet | Doppler (Minimal Wear)",
				Metadata: model.TransactionMetadata{
					AssetId: "456",
				},
				CreatedAt: time.Now(),
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
				PreviewUrl: "New Preview URL",
				CreatedAt:  time.Now(),
				Metadata: model.TransactionMetadata{
					AssetId: "123",
					Market:  "buff",
				},
			},
			{
				Name:      "★ Bayonet | Doppler (Minimal Wear)",
				CreatedAt: time.Now(),
				Metadata: model.TransactionMetadata{
					AssetId: "101112",
				},
			},
		}

		err = repo.UpsertTransactionsByAssetID(newTransactions)

		if err != nil {
			t.Error(err)
		}

		// Get the item back
		updatedTransaction, err := repo.FindTransactionByAssetId(transactions[0].Metadata.AssetId)
		if err != nil {
			t.Error(err)
		}

		// If dup key transaction already exists, the preview url should not be updated
		if updatedTransaction.PreviewUrl == newTransactions[0].PreviewUrl {
			t.Errorf("Preview URL SHALL NOT be updated: %v", updatedTransaction.PreviewUrl)
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
	db, repos, err := RepoInit()
	defer db.Disconnect()
	if err != nil {
		t.Error(err)
	}

	repo := repos.GetListingRepository()

	t.Run("DeleteOldListing", func(t *testing.T) {
		p1, _ := primitive.ParseDecimal128("100")
		p2, _ := primitive.ParseDecimal128("105")
		listings := []model.Listing{
			{
				Name:    "★ Bayonet | Doppler (Factory New)",
				AssetId: "123",
				Price:   p1,
			},
			{
				Name:    "★ Bayonet | Doppler (Factory New)",
				AssetId: "123",
				Price:   p2,
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
	db, repos, err := RepoInit()
	defer db.Disconnect()
	if err != nil {
		t.Error(err)
	}

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
			// transactionRepo.DeleteOldTransactionsByAssetID(assetId)
		}
	})
}

func TestSubscriptions(t *testing.T) {
	db, repos, err := RepoInit()
	defer db.Disconnect()
	if err != nil {
		t.Error(err)
	}

	repo := repos.GetSubscriptionRepository()

	t.Run("Subscriptions", func(t *testing.T) {
		subscriptions := model.Subscription{
			Name:       "★ Bayonet | Marble Fade (Factory New)",
			Rarity:     "FFI",
			MaxPremium: "5%",
			NotiType:   "telegram",
			NotiId:     "123",
		}

		_, err = repo.InsertSubscription(&subscriptions)

		if err != nil {
			t.Error(err)
		}

		// delete the subscription by name
		err = repo.DeleteSubscriptionByName(subscriptions.Name)
		if err != nil {
			t.Error(err)
		}

	})
}
