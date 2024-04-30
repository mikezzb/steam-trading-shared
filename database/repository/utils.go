package repository

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/mikezzb/steam-trading-shared/database/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GenerateUpdateBson(oldValue, newValue interface{}) bson.M {
	updateDoc := bson.M{}

	oldValueReflect := reflect.ValueOf(oldValue)
	newValueReflect := reflect.ValueOf(newValue)

	for i := 0; i < oldValueReflect.NumField(); i++ {
		oldField := oldValueReflect.Field(i)
		newField := newValueReflect.Field(i)

		if !reflect.DeepEqual(oldField.Interface(), newField.Interface()) {
			fieldName := oldValueReflect.Type().Field(i).Name
			updateDoc[fieldName] = newField.Interface()
		}
	}

	return updateDoc
}

// func SetUpdatedAtBson(m bson.M) bson.M {
// 	timestamp := GetTimestampNow()
// 	m["updatedAt"] = timestamp
// 	return m
// }

func GetBsonWithUpdatedAt() bson.M {
	timestamp := GetTimestampNow()
	return bson.M{"updatedAt": timestamp}
}

func GetTimestampNow() string {
	currTime := time.Now()
	timestamp := currTime.UnixNano() / int64(time.Millisecond)
	return strconv.FormatInt(timestamp, 10)
}

func MapToBson(m map[string]interface{}) bson.M {
	b := bson.M{}
	if m == nil {
		return b
	}

	for k, v := range m {
		b[k] = v
	}
	return b
}

func GetTransactionKey(tran *model.Transaction) string {
	return fmt.Sprintf("%s-%s", tran.Metadata.AssetId, tran.Metadata.Market)
}

// Page starts from 1
func GetPageOpts(page, pageSize int) *options.FindOptions {
	return options.Find().SetSkip(int64((page - 1) * pageSize)).SetLimit(int64(pageSize))
}

var ErrDuplicate = fmt.Errorf("duplicate key error")
