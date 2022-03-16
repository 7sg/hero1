package db

import (
	"hero1/go-services/historical-events/internal/config"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// NewClient create the new db client
func NewClient(conf *config.Config) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // client connection timeout
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.DB.Uri))
	return client, err
}
