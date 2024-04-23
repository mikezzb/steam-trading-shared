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

type TransactionRepository struct {
	TransactionCol       *mongo.Collection
	ChangeStreamCallback ChangeStreamCallback
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
		filter := bson.M{
			"assetId": transaction.AssetId,
			"market":  transaction.Market,
		}
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

func (r *TransactionRepository) DeleteOldTransactionsByAssetID(assetID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	// Find all transactions with the specified asset ID, sorted by insertion time in descending order
	cursor, err := r.TransactionCol.Find(ctx, bson.M{"assetId": assetID}, options.Find().SetSort(bson.D{{Key: "_id", Value: -1}}))
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	// Keep track of the latest transaction ID
	var latestTransactionID primitive.ObjectID
	if cursor.Next(ctx) {
		var latestTransaction model.Transaction
		if err := cursor.Decode(&latestTransaction); err != nil {
			return err
		}
		latestTransactionID = latestTransaction.ID
	}

	// Delete all transactions with the specified asset ID except the latest one
	_, err = r.TransactionCol.DeleteMany(ctx,
		bson.M{
			"assetId": assetID,
			"_id": bson.M{
				"$ne": latestTransactionID,
			},
		},
	)
	return err
}

func (r *TransactionRepository) GetAllUniqueAssetIDs() ([]string, error) {
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
	cursor, err := r.TransactionCol.Aggregate(ctx, pipeline)
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
