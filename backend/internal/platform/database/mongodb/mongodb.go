package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func createConnection(ctx context.Context, conn_str string) (*mongo.Client, error) {
	clientOptions := options.Client().
		ApplyURI(conn_str).
		SetMaxPoolSize(100).
		SetTimeout(time.Second * 10)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return client, nil
}

func MongoDB(ctx context.Context, conn_str string) (*mongo.Client, error) {
	for i := range 4 {
		fmt.Printf("Attempting to instrument mongo database connection - attempt#[%d]\n", i)
		mongo_client, err := createConnection(ctx, conn_str)
		if err != nil {
			fmt.Printf("Attempt failed due: %v\n", err)
			continue
		}

		return mongo_client, nil
	}

	return nil, fmt.Errorf("MongoDB connection attempts failed")

}
