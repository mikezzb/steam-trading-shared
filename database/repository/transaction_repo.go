package repository

import (
	"context"
	"fmt"

	"github.com/mikezzb/steam-trading-shared/database/model"
	"go.mongodb.org/mongo-driver/bson"
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
	err := r.TransactionCol.FindOne(ctx, bson.M{"metadata.assetId": assetId}).Decode(&transaction)
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

	var (
		operations  []mongo.WriteModel
		filters     []interface{}
		newTransMap = make(map[string]bool)
	)

	// Collect all filters and prepare to check in bulk
	for _, transaction := range transactions {
		filter := bson.M{
			"metadata.assetId": transaction.Metadata.AssetId,
			"metadata.market":  transaction.Metadata.Market,
		}
		filters = append(filters, filter)
		uniqueKey := GetTransactionKey(&transaction)
		newTransMap[uniqueKey] = true
	}

	// Perform a bulk check for existing documents
	if len(filters) > 0 {
		cursor, err := r.TransactionCol.Find(ctx, bson.M{"$or": filters})
		if err != nil {
			return err
		}
		defer cursor.Close(ctx)

		// Iterate through the cursor and mark existing documents
		var result bson.M
		for cursor.Next(ctx) {
			if err := cursor.Decode(&result); err != nil {
				continue
			}
			assetId := result["metadata"].(bson.M)["assetId"].(string)
			market := result["metadata"].(bson.M)["market"].(string)
			uniqueKey := fmt.Sprintf("%s-%s", assetId, market)
			delete(newTransMap, uniqueKey)
		}
	}

	// Prepare insert operations for non-existing documents
	for _, transaction := range transactions {
		uniqueKey := GetTransactionKey(&transaction)
		if _, found := newTransMap[uniqueKey]; found {
			model := mongo.NewInsertOneModel().SetDocument(transaction)
			operations = append(operations, model)
		}
	}

	// Execute bulk insert for new transactions
	if len(operations) > 0 {
		_, err := r.TransactionCol.BulkWrite(ctx, operations)
		if err != nil {
			return err
		}
	}
	return nil
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
