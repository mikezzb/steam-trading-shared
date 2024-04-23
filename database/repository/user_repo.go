package repository

import (
	"context"

	"github.com/mikezzb/steam-trading-shared/database/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	UserCol *mongo.Collection
}

// @return user, error (if error is nil but NO user, wrong password, otherwise username DNE)
func (r *UserRepository) GetUser(username, password string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_DURATION)
	defer cancel()

	user := &model.User{}
	filter := bson.M{"username": username}
	err := r.UserCol.FindOne(ctx, filter).Decode(user)

	if err != nil {
		return nil, err
	}

	if user.Password != password {
		return nil, nil
	}

	return user, err
}
