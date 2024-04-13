package repository

import (
	"context"

	"github.com/mikezzb/steam-trading-shared/database/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TransactionRepository struct {
	TransactionCol *mongo.Collection
}

func (r *TransactionRepository) FindTransactionByItemName(name string) (*model.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	var transaction model.Transaction
	err := r.TransactionCol.FindOne(ctx, bson.M{"name": name}).Decode(&transaction)
	return &transaction, err
}

func (r *TransactionRepository) InsertTransactions(transactions []model.Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()
	var documents []interface{}
	for _, transaction := range transactions {
		transaction.ID = primitive.NewObjectID()
		documents = append(documents, transaction)
	}
	_, err := r.TransactionCol.InsertMany(ctx, documents, options.InsertMany())
	return err
}

func (r *TransactionRepository) DeleteTransactionByItemName(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	_, err := r.TransactionCol.DeleteOne(ctx,
		bson.M{"name": name},
	)
	return err
}
