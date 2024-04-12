package repository

import (
	"context"
	"fmt"
	"steam-trading/shared/database/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ItemRepository struct {
	ItemCol *mongo.Collection
}

func (r *ItemRepository) FindItemByName(name string) (*model.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var item model.Item
	err := r.ItemCol.FindOne(ctx, bson.M{"name": name}).Decode(&item)
	return &item, err
}

func sameItem(item1, item2 *model.Item) bool {
	if item1 == nil || item2 == nil {
		return false
	}
	return item1.Name == item2.Name &&
		item1.IconUrl == item2.IconUrl &&
		item1.LowestMarketPrice == item2.LowestMarketPrice &&
		item1.LowestMarketName == item2.LowestMarketName &&
		item1.SteamPrice == item2.SteamPrice
}

func (r *ItemRepository) DeleteItemByName(item model.Item) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
	item := bson.M{}
	if oldItem.Name != newItem.Name {
		item["name"] = newItem.Name
	}
	if oldItem.IconUrl != newItem.IconUrl {
		item["icon_url"] = newItem.IconUrl
	}
	if oldItem.LowestMarketPrice != newItem.LowestMarketPrice {
		item["lowest_market_price"] = newItem.LowestMarketPrice
	}
	if oldItem.LowestMarketName != newItem.LowestMarketName {
		item["lowest_market_name"] = newItem.LowestMarketName
	}
	if oldItem.SteamPrice != newItem.SteamPrice {
		item["steam_price"] = newItem.SteamPrice
	}
	return bson.M{
		"$set": item,
	}
}

func (r *ItemRepository) UpdateItem(item model.Item) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// find existing item by name
	oldItem, err := r.FindItemByName(item.Name)
	// sync price if oldItem exists
	if err == nil {
		// if item is NOT better than oldItem, keep oldItem's price
		if item.Name != oldItem.Name && item.LowestMarketPrice > oldItem.LowestMarketPrice {
			item.LowestMarketPrice = oldItem.LowestMarketPrice
			item.LowestMarketName = oldItem.LowestMarketName
		}
	}
	// upsert if has delta
	if !sameItem(oldItem, &item) {
		opt := options.Update().SetUpsert(true)
		update := GetItemUpdateBson(oldItem, &item)
		fmt.Printf("bson: %v\n", GetItemUpdateBson(oldItem, &item))

		_, err := r.ItemCol.UpdateOne(ctx, bson.M{"name": item.Name}, update, opt)
		return err
	}
	return nil
}
