package repository

import (
	"context"

	"github.com/mikezzb/steam-trading-shared/database/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	UserCol *mongo.Collection
}

// @return user, error
func (r *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	user := &model.User{}
	filter := bson.M{"email": email}
	err := r.UserCol.FindOne(ctx, filter).Decode(user)

	if err != nil {
		return nil, err
	}

	return user, err
}

func (r *UserRepository) GetUserById(id primitive.ObjectID) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	user := &model.User{}
	err := r.UserCol.FindOne(ctx, bson.M{"_id": id}).Decode(user)
	return user, err
}

func (r *UserRepository) InsertUser(user *model.User) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	// ensure there are no dup in username OR email
	filter := bson.M{"$or": []bson.M{
		{"username": user.Username},
		{"email": user.Email},
	}}

	// check if user already exists
	count, err := r.UserCol.CountDocuments(ctx, filter)
	if err != nil {
		return primitive.NilObjectID, err
	}

	if count > 0 {
		return primitive.NilObjectID, ErrDuplicate
	}

	result, err := r.UserCol.InsertOne(ctx, user)
	return result.InsertedID.(primitive.ObjectID), err
}
