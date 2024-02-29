package common

import (
	"context"
	"errors"
	"fmt"
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
	uri2 := os.Getenv("REDIS_URI")

	if uri == "" {
		return errors.New("you must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	if uri2 == "" {
		return errors.New("you must set your 'REDIS_URI' environmental variable")

	}
	// client2, err2 := redis.NewClient(context.Background(),options.Client().ApplyURI(uri2))
	// if err2!=nil{
	// 	return err
	// }

	db = client.Database("go_demo")
	// db2=client2.Ping("go_demo")

	client2 := redis.NewClient(&redis.Options{
		Addr: uri2,
	})

	// Check for connection errors
	ctx := context.Background()
	err2 := client2.Ping(ctx).Err()
	if err2 != nil {
		fmt.Println("Error connecting to Redis:", err2)
		return err
	}

	err3 := client2.Set(ctx, "key", "value", 0).Err()
	if err3 != nil {
		fmt.Println("Error setting string:", err)
	}

	fmt.Println("Connected to Redis successfully!")

	return nil
}

func CloseDB() error {
	return db.Client().Disconnect(context.Background())
}
