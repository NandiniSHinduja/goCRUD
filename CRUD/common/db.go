package common

import (
	"context"
	"errors"

	// "fmt"
	"os"

	"github.com/go-redis/redis/v8"

	// "github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database
var db2 *redis.Client

func GetDBCollection(col string) *mongo.Collection {
	return db.Collection(col)
}

func InitDB() error {
	uri := os.Getenv("MONGODB_URI")

	if uri == "" {
		return errors.New("you must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	db = client.Database("go_demo")

	return nil
}

func CloseDB() error {
	return db.Client().Disconnect(context.Background())
}
