package database

import (
	"context"
	"log"
	"task_manager/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	Client   *mongo.Client
	Database *mongo.Database
}

var cfg *config.Config = config.LoadConfig()

func Connect() (*MongoClient, error) {

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ContextTimeout)
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
	err := mc.Client.Disconnect(ctx)
	if err != nil {
		log.Fatalf("Ошибка отключения от базы данных: %v", err)
	}
	log.Println("Отключение от базы данных успешно")
	return nil
}
