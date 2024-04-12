package repository

import "github.com/mikezzb/steam-trading-shared/database"

type Repositories struct {
	ItemRepo        ItemRepository
	ListingRepo     ListingRepository
	TransactionRepo TransactionRepository
	// NotiRepo        NotificationRepository
	// UserRepo        UserRepository
}

func NewRepositories(dbClient *database.DBClient) *Repositories {
	return &Repositories{
		ItemRepo:        *NewItemRepository(dbClient),
		ListingRepo:     *NewListingRepository(dbClient),
		TransactionRepo: *NewTransactionRepository(dbClient),
		// NotiRepo:        notiRepo,
		// UserRepo:        userRepo,
	}
}
