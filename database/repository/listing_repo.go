package repository

import (
	"context"
	"fmt"

	"github.com/mikezzb/steam-trading-shared/database/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ListingRepository struct {
	ListingCol           *mongo.Collection
	ChangeStreamCallback ChangeStreamCallback
}

func (r *ListingRepository) FindListingByItemName(name string) (*model.Listing, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	var listing model.Listing
	err := r.ListingCol.FindOne(ctx, bson.M{"name": name}).Decode(&listing)
	return &listing, err
}

func (r *ListingRepository) FindOneListing(filter bson.M) (*model.Listing, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	var listing model.Listing
	err := r.ListingCol.FindOne(ctx, filter).Decode(&listing)
	return &listing, err
}

func (r *ListingRepository) InsertListings(listings []model.Listing) error {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()
	var documents []interface{}
	for _, listing := range listings {
		documents = append(documents, listing)
	}
	_, err := r.ListingCol.InsertMany(ctx, documents, options.InsertMany())
	return err
}

// Returns the listings that were really updated / created
func (r *ListingRepository) UpsertListingsByAssetID(listings []model.Listing) ([]model.Listing, error) {
	updatedListings := make([]model.Listing, 0)

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	// upsert each listing to know which one was really updated / created
	for _, listing := range listings {
		filter := bson.M{"assetId": listing.AssetId}
		update := bson.M{"$set": listing}
		opts := options.Update().SetUpsert(true)

		result, err := r.ListingCol.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			return nil, err
		}

		// if created or updated
		if result.UpsertedID != nil || result.ModifiedCount > 0 {
			updatedListings = append(updatedListings, listing)
		}
	}

	return updatedListings, nil
}

func (r *ListingRepository) BulkUpsertListingsByAssetID(listings []model.Listing) error {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	var operations []mongo.WriteModel
	for _, listing := range listings {
		filter := bson.M{"assetId": listing.AssetId}
		update := bson.M{"$set": listing}
		model := mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true)
		operations = append(operations, model)
	}

	_, err := r.ListingCol.BulkWrite(ctx, operations)
	return err
}

func (r *ListingRepository) DeleteListingByItemName(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	_, err := r.ListingCol.DeleteOne(ctx,
		bson.M{"name": name},
	)
	return err
}

func (r *ListingRepository) FindItemByAssetId(assetID string) (*model.Listing, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	var listing model.Listing
	err := r.ListingCol.FindOne(ctx, bson.M{"assetId": assetID}).Decode(&listing)
	return &listing, err
}

func (r *ListingRepository) DeleteOldListingsByAssetID(assetID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	// Find all listings with the specified asset ID, sorted by insertion time in descending order
	cursor, err := r.ListingCol.Find(ctx, bson.M{"assetId": assetID}, options.Find().SetSort(bson.D{{Key: "_id", Value: -1}}))
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	// Keep track of the latest listing ID
	var latestListingID primitive.ObjectID
	if cursor.Next(ctx) {
		var latestListing model.Listing
		if err := cursor.Decode(&latestListing); err != nil {
			return err
		}
		latestListingID = latestListing.ID
	}

	// Delete all listings with the specified asset ID except the latest one
	_, err = r.ListingCol.DeleteMany(ctx,
		bson.M{
			"assetId": assetID,
			"_id": bson.M{
				"$ne": latestListingID,
			},
		},
	)
	return err
}

func (r *ListingRepository) GetAllUniqueAssetIDs() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	// Aggregate operation to get unique asset IDs
	pipeline := bson.A{
		bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$assetId"},
			}},
		},
	}

	// Execute aggregation
	cursor, err := r.ListingCol.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Iterate over cursor and collect unique asset IDs
	var assetIDs []string
	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		assetID, ok := result["_id"].(string)
		if !ok {
			return nil, fmt.Errorf("unexpected asset ID type")
		}
		assetIDs = append(assetIDs, assetID)
	}

	return assetIDs, nil
}

func (r *ListingRepository) DeleteAll() error {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	_, err := r.ListingCol.DeleteMany(ctx, bson.M{})
	return err
}
