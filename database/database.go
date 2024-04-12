package database

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	dbClientInstance *DBClient
	dbClientOnce     sync.Once
)

type DBClient struct {
	client *mongo.Client
	DB     *mongo.Database
}

func NewDBClient(uri, dbName string, timeout time.Duration) (*DBClient, error) {
	dbClientOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if err != nil {
			panic(err)
		}

		if err = client.Ping(ctx, nil); err != nil {
			client.Disconnect(ctx)
			panic(err)
		}

		db := client.Database(dbName)
		dbClientInstance = &DBClient{
			client: client,
			DB:     db,
		}
	})

	return dbClientInstance, nil
}

func (c *DBClient) Disconnect() error {
	return c.client.Disconnect(context.Background())
}

func (c *DBClient) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return c.client.Ping(ctx, nil)
}
