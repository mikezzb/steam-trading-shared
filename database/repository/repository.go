package repository

import "github.com/mikezzb/steam-trading-shared/database"

type RepoFactory interface {
	GetItemRepository() *ItemRepository
	GetListingRepository() *ListingRepository
	GetTransactionRepository() *TransactionRepository
	GetSubscriptionRepository() *SubscriptionRepository
}

type Repositories struct {
	dbClient             *database.DBClient
	changeStreamHandlers *ChangeStreamHandlers
	itemRepo             *ItemRepository
	listingRepo          *ListingRepository
	transactionRepo      *TransactionRepository
	subscriptionRepo     *SubscriptionRepository
}

type ChangeStreamHandlers struct {
	ItemChangeStreamCallback         ChangeStreamCallback
	ListingChangeStreamCallback      ChangeStreamCallback
	TransactionChangeStreamCallback  ChangeStreamCallback
	SubscriptionChangeStreamCallback ChangeStreamCallback
}

type ChangeStreamCallback func(data interface{}, operationType string)

func NewRepoFactory(dbClient *database.DBClient, handlers *ChangeStreamHandlers) *Repositories {
	if handlers == nil {
		handlers = &ChangeStreamHandlers{}
	}
	return &Repositories{
		dbClient:             dbClient,
		changeStreamHandlers: handlers,
	}
}

// factory
func (r *Repositories) GetItemRepository() *ItemRepository {
	if r.itemRepo == nil {
		r.itemRepo = &ItemRepository{
			ItemCol:              r.dbClient.DB.Collection("items"),
			ChangeStreamCallback: r.changeStreamHandlers.ItemChangeStreamCallback,
		}
	}
	return r.itemRepo
}

func (r *Repositories) GetListingRepository() *ListingRepository {
	if r.listingRepo == nil {
		r.listingRepo = &ListingRepository{
			ListingCol:           r.dbClient.DB.Collection("listings"),
			ChangeStreamCallback: r.changeStreamHandlers.ListingChangeStreamCallback,
		}
	}
	return r.listingRepo
}

func (r *Repositories) GetTransactionRepository() *TransactionRepository {
	if r.transactionRepo == nil {
		r.transactionRepo = &TransactionRepository{
			TransactionCol:       r.dbClient.DB.Collection("transactions"),
			ChangeStreamCallback: r.changeStreamHandlers.TransactionChangeStreamCallback,
		}
	}
	return r.transactionRepo
}

func (r *Repositories) GetSubscriptionRepository() *SubscriptionRepository {
	if r.subscriptionRepo == nil {
		r.subscriptionRepo = &SubscriptionRepository{
			SubCol:               r.dbClient.DB.Collection("subscriptions"),
			ChangeStreamCallback: r.changeStreamHandlers.SubscriptionChangeStreamCallback,
		}
	}
	return r.subscriptionRepo
}
