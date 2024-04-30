package database

import (
	"context"
	"log"
	"strings"
	"time"

	shared "github.com/mikezzb/steam-trading-shared"
	"github.com/mikezzb/steam-trading-shared/database/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// call this function to init collections, only call once
func (c *DBClient) Init() error {
	// init timeseries collection
	tsCollName := "transactions"

	tsOpts := options.TimeSeries().SetTimeField("createdAt").SetMetaField("metadata").SetGranularity("hours")

	createCollOpts := options.CreateCollection().SetTimeSeriesOptions(tsOpts)

	if err := c.DB.CreateCollection(context.Background(), tsCollName, createCollOpts); err != nil {
		return err
	}

	return nil
}

// migrate
func (c *DBClient) MigrateTransactions(oldCollName, newCollName string) error {
	oldColl := c.DB.Collection(oldCollName)
	newColl := c.DB.Collection(newCollName)

	cursor, err := oldColl.Find(context.Background(), bson.D{})
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return err
		}

		delete(doc, "updatedAt")

		// print type of createdAt
		log.Printf("Type of createdAt: %T", doc["createdAt"])

		createdAt, err := unixToTime(doc["createdAt"])
		if err != nil {
			return err
		}

		doc["createdAt"] = createdAt

		if _, err := newColl.InsertOne(context.Background(), doc); err != nil {
			return err
		}
	}

	if err := cursor.Err(); err != nil {
		return err
	}

	log.Println("Data migration completed.")

	return nil
}

func unixToTime(createdAt interface{}) (time.Time, error) {
	if createdAt32, ok := createdAt.(int32); ok {
		return time.Unix(int64(createdAt32), 0), nil
	} else if createdAt64, ok := createdAt.(int64); ok {
		return time.Unix(createdAt64, 0), nil
	}

	return time.Time{}, nil
}

func (c *DBClient) ConvertUnixToTime(coll string) error {
	// convert unix timestamp to time.Time
	cursor, err := c.DB.Collection(coll).Find(context.Background(), bson.D{})
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return err
		}

		if _, ok := doc["createdAt"]; !ok {
			if createdAt, err := unixToTime(doc["timestamp"]); err == nil {
				doc["createdAt"] = createdAt
			} else {
				return err
			}
		}

		if _, ok := doc["updatedAt"]; !ok {
			if updatedAt, err := unixToTime(doc["timestamp"]); err == nil {
				doc["updatedAt"] = updatedAt
			} else {
				return err
			}
		}

		if _, err := c.DB.Collection(coll).UpdateOne(context.Background(), bson.M{"_id": doc["_id"]}, bson.M{"$set": doc}); err != nil {
			return err
		}
	}

	if err := cursor.Err(); err != nil {
		return err
	}

	log.Println("Data conversion completed.")

	return nil
}

// convert string to Decimal128
func val2decimal128(val interface{}) (primitive.Decimal128, error) {
	if val == nil {
		return primitive.ParseDecimal128("error")
	}
	if valStr, ok := val.(string); ok {
		return primitive.ParseDecimal128(valStr)
	}

	return primitive.ParseDecimal128("error")
}

// recursively convert all specified fields to Decimal128
func (c *DBClient) ConvertToDecimal128(collName string, fields []string) error {
	cursor, err := c.DB.Collection(collName).Find(context.Background(), bson.D{})
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return err
		}

		for _, field := range fields {
			nestedFields := strings.Split(field, ".")
			nestedDoc := doc
			for _, nestedField := range nestedFields[:len(nestedFields)-1] {
				if nestedVal, ok := nestedDoc[nestedField]; ok {
					nestedDoc = nestedVal.(bson.M)
				} else {
					log.Printf("Nested Field %v does not exist in document", nestedField)
					continue
				}
			}
			priceField := nestedFields[len(nestedFields)-1]
			if val, ok := nestedDoc[priceField]; ok {
				val, err := val2decimal128(val)
				if err != nil {
					log.Printf("Failed to convert field %s value %v to Decimal128: %v", field, val, err)
					delete(nestedDoc, priceField)
				} else {
					nestedDoc[priceField] = val
				}
			} else {
				log.Printf("Field %v does not exist in document", field)
				continue
			}
		}

		if _, err := c.DB.Collection(collName).UpdateOne(context.Background(), bson.M{"_id": doc["_id"]}, bson.M{"$set": doc}); err != nil {
			return err
		}
	}

	if err := cursor.Err(); err != nil {
		return err
	}

	log.Println("Data conversion completed.")
	return nil
}

// Reformat the transaction collection
func (c *DBClient) ReformatTransactionCollection(collName string) error {
	cursor, err := c.DB.Collection(collName).Find(context.Background(), bson.D{})
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return err
		}

		// make metadata field
		market, ok := doc["market"].(string)
		if !ok {
			log.Printf("Set market to buff by default")
			market = "buff"
		}

		matadata := model.TransactionMetadata{
			AssetId: doc["assetId"].(string),
			Market:  market,
		}

		delete(doc, "assetId")
		delete(doc, "market")

		doc["metadata"] = matadata

		log.Printf("Reformatted document: %v", doc)

		// update the document
		if _, err := c.DB.Collection(collName).ReplaceOne(context.Background(), bson.M{"_id": doc["_id"]}, doc); err != nil {
			return err
		}
	}

	if err := cursor.Err(); err != nil {
		return err
	}

	log.Println("Data reformatting completed.")
	return nil
}

func (c *DBClient) ReformatItems(collName string) error {
	cursor, err := c.DB.Collection(collName).Find(context.Background(), bson.D{})
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return err
		}

		// augment the name
		category, skin, exterior := shared.DecodeItemFullName(doc["name"].(string))

		doc["category"] = category
		doc["skin"] = skin
		doc["exterior"] = exterior

		log.Printf("Category: %v, Skin: %v, Exterior: %v", category, skin, exterior)

		// update the document
		if _, err := c.DB.Collection(collName).ReplaceOne(context.Background(), bson.M{"_id": doc["_id"]}, doc); err != nil {
			return err
		}
	}

	if err := cursor.Err(); err != nil {
		return err
	}

	log.Println("Data reformatting completed.")
	return nil
}
