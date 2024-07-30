package db

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	URI    string
	DBName string
}

func New(URI, DBName string) *MongoConfig {
	return &MongoConfig{
		URI:    URI,
		DBName: DBName,
	}
}

func (cfg *MongoConfig) Connect() (*mongo.Database, error) {
	if cfg.URI == "" {
		return nil, errors.New("you must set your 'uri' environmental variable")
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, err
	}

	return client.Database(cfg.DBName), nil
}

func (cfg *MongoConfig) Close(db *mongo.Database) error {
	return db.Client().Disconnect(context.Background())
}
