package app

import (
	"context"
	"fmt"
	"time"

	"github.com/ronnyp07/SportStream/worker/internal/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	defaultTimeout = 10 * time.Second
)

type MongoDB struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func SetupMongoDB(ctx context.Context) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	// Get configuration
	cfg := config.Infra().MongoDB
	if cfg.Url == "" {
		return nil, fmt.Errorf("MongoDB URL configuration is missing")
	}

	clientOptions := options.Client().
		ApplyURI(cfg.Url).
		SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1)).
		SetAuth(options.Credential{
			AuthSource:    cfg.AuthSource,
			Username:      cfg.UserName,
			Password:      cfg.PassWord,
			AuthMechanism: "SCRAM-SHA-256",
		}).
		SetConnectTimeout(defaultTimeout).
		SetSocketTimeout(30 * time.Second)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := client.Database(cfg.DataBase)

	return &MongoDB{
		Client: client,
		DB:     db,
	}, nil
}

func (m *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	return m.Client.Disconnect(ctx)
}
