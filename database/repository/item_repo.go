package repository

import (
	"context"

	"github.com/mikezzb/steam-trading-shared/database/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ItemRepository struct {
	ItemCol *mongo.Collection
}

func (r *ItemRepository) FindItemByName(name string) (*model.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
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
	if oldItem.IconUrl != newItem.IconUrl {
		item["iconUrl"] = newItem.IconUrl
	}
	if oldItem.LowestMarketPrice != newItem.LowestMarketPrice {
		item["lowestMarketPrice"] = newItem.LowestMarketPrice
	}
	if oldItem.LowestMarketName != newItem.LowestMarketName {
		item["lowestMarketName"] = newItem.LowestMarketName
	}
	if oldItem.SteamPrice != newItem.SteamPrice {
		item["steamPrice"] = newItem.SteamPrice
	}
	return bson.M{
		"$set": item,
	}
}

func (r *ItemRepository) UpdateItem(item *model.Item) error {
	// clone item for modification
	item = &model.Item{
		Name:              item.Name,
		IconUrl:           item.IconUrl,
		LowestMarketPrice: item.LowestMarketPrice,
		LowestMarketName:  item.LowestMarketName,
		SteamPrice:        item.SteamPrice,
	}
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
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
	if !sameItem(oldItem, item) {
		opt := options.Update().SetUpsert(true)
		update := GetItemUpdateBson(oldItem, item)

		_, err := r.ItemCol.UpdateOne(ctx, bson.M{"name": item.Name}, update, opt)
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
