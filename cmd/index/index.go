package main

import (
	"context"
	"optimistic-booking-system/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func createUniqueIndex(collection *mongo.Collection) error {
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
			{Key: "resource_id", Value: 1},
			{Key: "start_time", Value: 1},
			{Key: "end_time", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	_, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
	return err
}

func main() {
	mongoCfg := db.New("mongodb://root:example@localhost:27017", "optimistic")
	mongo, err := mongoCfg.Connect()
	if err != nil {
		panic(err)
	}
	defer mongoCfg.Close(mongo)

	bookingCollection := mongo.Collection("bookings")
	err = createUniqueIndex(bookingCollection)
	if err != nil {
		panic(err)
	}
}
