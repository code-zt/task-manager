package database

import (
	"context"
	"log"
	"task_manager/internal/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	Client   *mongo.Client
	Database *mongo.Database
}

var cfg *config.Config

func Connect() (*MongoClient, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфига: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+cfg.DatabaseHost+":"+cfg.DatabasePort))
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &MongoClient{
		Client:   client,
		Database: client.Database(cfg.DatabaseName),
	}, nil
}
func (mc *MongoClient) Disconnect(ctx context.Context) error {
	return mc.Client.Disconnect(ctx)
}
