package repository

import (
	"context"
	"log"
	"time"

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

func sameItem(item1, item2 *model.Item) bool {
	if item1 == nil || item2 == nil {
		return false
	}
	return item1.Name == item2.Name &&
		item1.BuffPrice.UpdatedAt == item2.BuffPrice.UpdatedAt &&
		item1.IgxePrice.UpdatedAt == item2.IgxePrice.UpdatedAt &&
		item1.UUPrice.UpdatedAt == item2.UUPrice.UpdatedAt &&
		item1.SteamPrice.UpdatedAt == item2.SteamPrice.UpdatedAt
}

func (r *ItemRepository) DeleteItemByName(item *model.Item) error {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	_, err := r.ItemCol.DeleteOne(ctx,
		bson.M{"name": item.Name},
	)
	return err
}

func GetItemUpdateBson(oldItem, newItem *model.Item) interface{} {
	if oldItem == nil {
		return bson.M{
			"$set": newItem,
		}
	}
	// find dirty fields
	item := GetBsonWithUpdatedAt()
	if oldItem.Name != newItem.Name {
		item["name"] = newItem.Name
	}

	if oldItem.BuffPrice.UpdatedAt != newItem.BuffPrice.UpdatedAt {
		item["buffPrice"] = newItem.BuffPrice
	}

	if oldItem.IgxePrice.UpdatedAt != newItem.IgxePrice.UpdatedAt {
		item["igxePrice"] = newItem.IgxePrice
	}

	if oldItem.UUPrice.UpdatedAt != newItem.UUPrice.UpdatedAt {
		item["uuPrice"] = newItem.UUPrice
	}

	if oldItem.SteamPrice.UpdatedAt != newItem.SteamPrice.UpdatedAt {
		item["steamPrice"] = newItem.SteamPrice
	}

	return bson.M{
		"$set": item,
	}
}

// Upsert item by id
func (r *ItemRepository) UpsertItem(item *model.Item) error {
	// clone item for modification
	item = &model.Item{
		ID:         item.ID,
		Name:       item.Name,
		SteamPrice: item.SteamPrice,
		BuffPrice:  item.BuffPrice,
		IgxePrice:  item.IgxePrice,
		UUPrice:    item.UUPrice,
	}
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	// find existing item by name
	oldItem, _ := r.FindItemById(item.ID)
	// upsert if has delta
	if !sameItem(oldItem, item) {
		opt := options.Update().SetUpsert(true)
		update := GetItemUpdateBson(oldItem, item)

		_, err := r.ItemCol.UpdateOne(ctx, bson.M{"_id": item.ID}, update, opt)
		return err
	}
	return nil
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

func (r *ItemRepository) GetItemsByPage(page, size int, filters map[string]interface{}) ([]model.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	log.Printf("timeout: %v", 1*time.Second)
	defer cancel()

	filtersBson := MapToBson(filters)

	opts := GetPageOpts(page, size)

	log.Printf("Executing MongoDB Find with filters: %v and opts: %v", filtersBson, opts)

	cursor, err := r.ItemCol.Find(ctx, filtersBson, opts)
	if err != nil {
		return nil, err
	}

	var items []model.Item
	err = cursor.All(ctx, &items)
	return items, err
}

func (r *ItemRepository) Count() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	return r.ItemCol.CountDocuments(ctx, bson.M{})
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
