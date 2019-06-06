package core

import (
	"context"
	"fmt"

	"github.com/mongodb/mongo-go-driver/mongo"
)

type key string

// DbURL context key
const (
	DbURL = key("DbURL")
)

// ConfigDB configures database with new client
func ConfigDB(ctx context.Context, db string) (*mongo.Database, error) {
	uri := fmt.Sprintf(`%s`, ctx.Value(DbURL))
	client, err := mongo.NewClient(uri)
	if err != nil {
		return nil, fmt.Errorf("couldn't connect to mongo: %v", err)
	}
	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("mongo client couldn't connect with background context: %v", err)
	}
	DB := client.Database(db)
	return DB, nil
}
