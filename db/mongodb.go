package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectMongoDB connects into MongoDB by given URI and database name.
func ConnectMongoDB(ctx context.Context, uri, dbname string) (*mongo.Client, *mongo.Database, error) {
	dbClient, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, nil, err
	}

	if err = dbClient.Ping(ctx, nil); err != nil {
		return dbClient, nil, err
	}

	return dbClient, dbClient.Database(dbname), nil
}
