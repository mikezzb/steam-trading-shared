package database

import "github.com/mikezzb/steam-trading-shared/database/repository"

type Repositories struct {
	dbClient        *DBClient
	itemRepo        *repository.ItemRepository
	listingRepo     *repository.ListingRepository
	transactionRepo *repository.TransactionRepository
}

func NewRepositories(dbClient *DBClient) *Repositories {
	return &Repositories{
		dbClient: dbClient,
	}
}

// factory
func (r *Repositories) GetItemRepository() *repository.ItemRepository {
	if r.itemRepo == nil {
		r.itemRepo = &repository.ItemRepository{
			ItemCol: r.dbClient.DB.Collection("items"),
		}
	}
	return r.itemRepo
}

func (r *Repositories) GetListingRepository() *repository.ListingRepository {
	if r.listingRepo == nil {
		r.listingRepo = &repository.ListingRepository{
			ListingCol: r.dbClient.DB.Collection("listings"),
		}
	}
	return r.listingRepo
}

func (r *Repositories) GetTransactionRepository() *repository.TransactionRepository {
	if r.transactionRepo == nil {
		r.transactionRepo = &repository.TransactionRepository{
			TransactionCol: r.dbClient.DB.Collection("transactions"),
		}
	}
	return r.transactionRepo
}
