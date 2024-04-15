package repository

import (
	"context"

	"github.com/mikezzb/steam-trading-shared/database/model"
	"go.mongodb.org/mongo-driver/bson"
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

func (r *TransactionRepository) FindTransactionByAssetId(assetId string) (*model.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	var transaction model.Transaction
	err := r.TransactionCol.FindOne(ctx, bson.M{"assetId": assetId}).Decode(&transaction)
	return &transaction, err
}

func (r *TransactionRepository) InsertTransactions(transactions []model.Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()
	var documents []interface{}
	for _, transaction := range transactions {
		documents = append(documents, transaction)
	}
	_, err := r.TransactionCol.InsertMany(ctx, documents, options.InsertMany())
	return err
}

func (r *TransactionRepository) UpsertTransactionsByAssetID(transactions []model.Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	var operations []mongo.WriteModel
	for _, transaction := range transactions {
		filter := bson.M{"assetId": transaction.AssetId}
		update := bson.M{"$set": transaction}
		model := mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true)
		operations = append(operations, model)
	}

	_, err := r.TransactionCol.BulkWrite(ctx, operations)
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

func (r *TransactionRepository) DeleteAll() error {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	_, err := r.TransactionCol.DeleteMany(ctx, bson.M{})
	return err
}

// get all transactions
func (r *TransactionRepository) FindAllTransactions() ([]model.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	cursor, err := r.TransactionCol.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []model.Transaction
	for cursor.Next(ctx) {
		var transaction model.Transaction
		if err := cursor.Decode(&transaction); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}
