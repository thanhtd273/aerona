package configs

import (
	"context"
	"fmt"
	"os"
	"time"

	"aerona.thanhtd.com/notification-service/internal/logging"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func NewMongoClient() *mongo.Client {
	logger, err := logging.NewLogger(os.Getenv("LOG_PATH"), os.Getenv("LOG_LEVEL"))

	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	uri := fmt.Sprintf("mongodb://%s:%s@%s/%s?authSource=admin", os.Getenv("MONGO_USERNAME"), os.Getenv("MONGO_PASSWORD"), os.Getenv("MONGO_HOST"), "aerona_ticket_database")
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to connect to MongoDB, error: %v", err))
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Fatal("Failed to ping MongoDB, error: %v", zap.Error(err))
	}
	logger.Info("Connect to MongoDB successfully!")

	return client
}

func GetCollection(client *mongo.Client, dbName, collectionName string) *mongo.Collection {
	return client.Database(dbName).Collection(collectionName)
}
