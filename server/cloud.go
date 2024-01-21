package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	connectionStringEnvVar = "AZURE_COSMOS_JARVIS_DATABASE_CONNECTION_CONNECTIONSTRING"
)

func getMongoCollectionConnection(collectionName string) (*mongo.Collection, context.Context, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	if connectionString, ok := os.LookupEnv(connectionStringEnvVar); ok {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
		if err != nil {
			log.Fatalf("Failed to connect to mongo: %s\n", err)
		}
		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			log.Fatalf("Failed to ping mongos: %s\n", err)
		}

		cleanup := func() {
			cancel()
			if err = client.Disconnect(ctx); err != nil {
				log.Printf("WARN: Failed to disconnect from Cosmos: %s\n", err)
			}
		}

		// If in development mode, append -dev to collection
		if gin.Mode() != gin.ReleaseMode {
			collectionName = fmt.Sprintf("%s-dev", collectionName)
		}

		collection := client.Database("db").Collection(collectionName)
		return collection, ctx, cleanup
	}

	// Fatal
	log.Fatalf("Failed to connect to mongo: '%s' not set \n", connectionStringEnvVar)
	return nil, nil, nil
}
