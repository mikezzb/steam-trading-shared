package repository

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/mikezzb/steam-trading-shared/database/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetUpdateBson(oldValue, newValue interface{}) (bson.M, error) {
	updateDoc := bson.M{}

	oldValueReflect := reflect.Indirect(reflect.ValueOf(oldValue))
	newValueReflect := reflect.Indirect(reflect.ValueOf(newValue))

	if oldValueReflect.Kind() != reflect.Struct || newValueReflect.Kind() != reflect.Struct {
		return nil, fmt.Errorf("both oldValue and newValue must be structs or pointers to structs")
	}

	for i := 0; i < oldValueReflect.NumField(); i++ {
		oldField := oldValueReflect.Field(i)
		newField := newValueReflect.Field(i)
		typeField := oldValueReflect.Type().Field(i)

		// Using struct tags to determine BSON field names
		bsonFieldName := typeField.Tag.Get("bson")
		if bsonFieldName == "" {
			bsonFieldName = typeField.Name // Fallback to Go field name if no BSON tag is present
		}

		// Check if field names contain ",omitempty" and remove it for update operation context
		if idx := strings.Index(bsonFieldName, ",omitempty"); idx != -1 {
			bsonFieldName = bsonFieldName[:idx]
		}

		// Only include fields that have changed
		if !reflect.DeepEqual(oldField.Interface(), newField.Interface()) {
			updateDoc[bsonFieldName] = newField.Interface()
		}
	}

	return updateDoc, nil
}

func GetUpsertBson(oldValue, newValue interface{}) (bson.M, error) {
	updateDoc := bson.M{}

	oldValueReflect := reflect.Indirect(reflect.ValueOf(oldValue))
	newValueReflect := reflect.Indirect(reflect.ValueOf(newValue))

	if oldValueReflect.Kind() != reflect.Struct || newValueReflect.Kind() != reflect.Struct {
		return nil, fmt.Errorf("both oldValue and newValue must be structs or pointers to structs")
	}

	// Create a map of the old values for quick lookup
	oldValuesMap := make(map[string]reflect.Value)
	for i := 0; i < oldValueReflect.NumField(); i++ {
		fieldName := oldValueReflect.Type().Field(i).Name
		oldValuesMap[fieldName] = oldValueReflect.Field(i)
	}

	for i := 0; i < newValueReflect.NumField(); i++ {
		newField := newValueReflect.Field(i)
		typeField := newValueReflect.Type().Field(i)

		// Using struct tags to determine BSON field names
		bsonFieldName := typeField.Tag.Get("bson")
		if bsonFieldName == "" {
			bsonFieldName = typeField.Name // Fallback to Go field name if no BSON tag is present
		}

		// Check if field names contain ",omitempty" and remove it for update operation context
		if idx := strings.Index(bsonFieldName, ",omitempty"); idx != -1 {
			bsonFieldName = bsonFieldName[:idx]

			// If the new value is empty, skip this field
			if reflect.DeepEqual(newField.Interface(), reflect.Zero(newField.Type()).Interface()) {
				continue
			}
		}

		// Check if the old value has this field; if it does, compare them
		if oldField, ok := oldValuesMap[typeField.Name]; ok {
			if !reflect.DeepEqual(oldField.Interface(), newField.Interface()) {
				updateDoc[bsonFieldName] = newField.Interface()
			}
		} else {
			// If the old value does not have this field, just set the new value
			updateDoc[bsonFieldName] = newField.Interface()
		}
	}

	return updateDoc, nil
}

func AddUpdatedAtToBson(b bson.M) {
	b["updatedAt"] = time.Now()
}

func GetBsonWithUpdatedAt() bson.M {
	return bson.M{"updatedAt": time.Now()}
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
