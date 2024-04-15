package repository

import (
	"context"

	"github.com/mikezzb/steam-trading-shared/database/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ListingRepository struct {
	ListingCol *mongo.Collection
}

func (r *ListingRepository) FindListingByItemName(name string) (*model.Listing, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	var listing model.Listing
	err := r.ListingCol.FindOne(ctx, bson.M{"name": name}).Decode(&listing)
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

func (r *ListingRepository) UpsertListingsByAssetID(listings []model.Listing) error {
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
