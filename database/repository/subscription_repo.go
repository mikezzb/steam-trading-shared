package repository

import (
	"context"

	"github.com/mikezzb/steam-trading-shared/database/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SubscriptionRepository struct {
	SubCol               *mongo.Collection
	ChangeStreamCallback ChangeStreamCallback
}

func (r *SubscriptionRepository) UpsertSubscription(subscription *model.Subscription) error {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()
	// upsert on name
	filter := bson.M{"name": subscription.Name}
	update := bson.M{"$set": subscription}
	opts := options.Update().SetUpsert(true)

	result, err := r.SubCol.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	if r.ChangeStreamCallback != nil {
		if result.UpsertedID != nil {
			r.ChangeStreamCallback(subscription, "insert")
		} else {
			r.ChangeStreamCallback(subscription, "update")
		}
	}

	return err
}

func (r *SubscriptionRepository) GetAll() ([]model.Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	cursor, err := r.SubCol.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var subscriptions []model.Subscription
	err = cursor.All(ctx, &subscriptions)
	return subscriptions, err
}

func (r *SubscriptionRepository) DeleteSubscriptionByName(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	result, err := r.SubCol.DeleteOne(ctx, bson.M{"name": name})

	if err != nil {
		return err
	}

	if r.ChangeStreamCallback != nil {
		r.ChangeStreamCallback(result, "delete")
	}

	return err
}

func (r *SubscriptionRepository) DeleteAll() error {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	_, err := r.SubCol.DeleteMany(ctx, bson.M{})
	return err
}
