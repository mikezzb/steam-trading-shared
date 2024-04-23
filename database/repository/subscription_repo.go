package repository

import (
	"context"

	"github.com/mikezzb/steam-trading-shared/database/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SubscriptionRepository struct {
	SubCol               *mongo.Collection
	ChangeStreamCallback ChangeStreamCallback
}

func (r *SubscriptionRepository) InsertSubscription(subscription *model.Subscription) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	result, err := r.SubCol.InsertOne(ctx, subscription)
	if err != nil {
		return primitive.NilObjectID, err
	}

	if r.ChangeStreamCallback != nil {
		r.ChangeStreamCallback(subscription, "insert")
	}

	return result.InsertedID.(primitive.ObjectID), err
}

func (r *SubscriptionRepository) UpdateSubscription(subscription *model.Subscription) error {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	_, err := r.SubCol.ReplaceOne(ctx, bson.M{"_id": subscription.ID}, subscription)
	if err != nil {
		return err
	}

	if r.ChangeStreamCallback != nil {
		r.ChangeStreamCallback(subscription, "update")
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
