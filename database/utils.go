package database

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
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
