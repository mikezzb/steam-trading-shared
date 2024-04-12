package repository

import (
	"github.com/mikezzb/steam-trading-shared/database"
)

type ListingRepository struct {
	dbClient *database.DBClient
}

func NewListingRepository(dbClient *database.DBClient) *ListingRepository {
	return &ListingRepository{
		dbClient: dbClient,
	}
}
