package repository

import (
	"context"
	"log"
	"time"

	shared "github.com/mikezzb/steam-trading-shared"
	"github.com/mikezzb/steam-trading-shared/database/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ItemRepository struct {
	ItemCol              *mongo.Collection
	ChangeStreamCallback ChangeStreamCallback
}

func (r *ItemRepository) FindItemByName(name string) (*model.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	var item model.Item
	err := r.ItemCol.FindOne(ctx, bson.M{"name": name}).Decode(&item)
	return &item, err
}

func (r *ItemRepository) FindItemById(id string) (*model.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	var item model.Item
	err := r.ItemCol.FindOne(ctx, bson.M{"_id": id}).Decode(&item)
	return &item, err
}

func (r *ItemRepository) DeleteItemByName(item *model.Item) error {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	_, err := r.ItemCol.DeleteOne(ctx,
		bson.M{"name": item.Name},
	)
	return err
}

// Upsert item by id
func (r *ItemRepository) UpsertItem(item *model.Item) error {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	// find existing item by name
	oldItem, _ := r.FindItemById(item.ID)
	// get upsert bson
	itemDelta, err := GetUpsertBson(oldItem, item)
	if err != nil {
		return err
	}
	AddUpdatedAtToBson(itemDelta)
	// if no change, return (but it shall not happen cuz the updatedAt field is always updated)
	if len(itemDelta) == 0 {
		return nil
	}
	update := bson.M{"$set": itemDelta}

	opt := options.Update().SetUpsert(true)

	_, err = r.ItemCol.UpdateOne(ctx, bson.M{"_id": item.ID}, update, opt)
	return err
}

func (r *ItemRepository) GetAll() ([]model.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	cursor, err := r.ItemCol.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var items []model.Item
	err = cursor.All(ctx, &items)
	return items, err
}

func (r *ItemRepository) GetItemsByPage(page, size int, filters bson.M) ([]model.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	log.Printf("timeout: %v", 1*time.Second)
	defer cancel()

	opts := GetPageOpts(page, size)

	log.Printf("Executing MongoDB Find with filters: %v and opts: %v", filters, opts)

	cursor, err := r.ItemCol.Find(ctx, filters, opts)
	if err != nil {
		return nil, err
	}

	var items []model.Item
	err = cursor.All(ctx, &items)
	return items, err
}

func (r *ItemRepository) Count(filters bson.M) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	return r.ItemCol.CountDocuments(ctx, filters)
}

func (r *ItemRepository) GetItemByName(name string) (*model.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	var item model.Item
	err := r.ItemCol.FindOne(ctx, bson.M{"name": name}).Decode(&item)
	return &item, err
}

func (r *ItemRepository) DeleteAll() error {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	_, err := r.ItemCol.DeleteMany(ctx, bson.M{})
	return err
}

// TODO: cache this in another collection using trigger
func (r *ItemRepository) GetItemFilters() (map[string][]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	// From all items, get all unique values of selected fields
	selectedFields := shared.ITEM_FIXED_VAL_FILTER_KEYS
	filters := make(map[string][]interface{})
	for _, field := range selectedFields {
		values, err := r.ItemCol.Distinct(ctx, field, bson.M{})
		if err != nil {
			return nil, err
		}
		filters[field] = values
	}

	return filters, nil
}
