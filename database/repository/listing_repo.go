package repository

import (
	"context"

	"github.com/mikezzb/steam-trading-shared/database/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
