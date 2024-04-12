package repository

import (
	"github.com/mikezzb/steam-trading-shared/database"
)

type TransactionRepository struct {
	dbClient *database.DBClient
}

func NewTransactionRepository(dbClient *database.DBClient) *TransactionRepository {
	return &TransactionRepository{
		dbClient: dbClient,
	}
}
