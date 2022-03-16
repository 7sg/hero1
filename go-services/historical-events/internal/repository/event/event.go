package event

import (
	"hero1/go-services/historical-events/internal/config"
	"hero1/go-services/historical-events/internal/domain/database"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

// Repository for events
type Repository struct {
	dbClient   *mongo.Client
	database   string
	collection string
	dbConfig   config.Db
}

// New will Create and instance of event.Repository
func New(mongoClient *mongo.Client, dbConfig config.Db) *Repository {
	return &Repository{
		dbClient:   mongoClient,
		database:   "hero1",
		collection: "event",
		dbConfig:   dbConfig,
	}
}

// Save persist the event in db
func (r *Repository) Save(ctx context.Context, event *database.Event) error {
	collection := r.dbClient.Database(r.database).Collection(r.collection)

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(r.dbConfig.WriteTimeOut)) // query execution timeout
	defer cancel()

	_, err := collection.InsertOne(ctx, event)
	return err
}

// Get find the event from db
func (r *Repository) Get(ctx context.Context, searchFilter *database.SearchFilter) ([]*database.Event, error) {
	collection := r.dbClient.Database(r.database).Collection(r.collection)

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(r.dbConfig.ReadTimeOut)) // query execution timeout
	defer cancel()

	cursor, err := collection.Find(ctx, searchFilter)
	if err != nil {
		log.Printf("error in Repository, from Get err:%+v", err)
		return nil, err
	}
	var events []*database.Event
	err = cursor.All(ctx, &events)

	log.Printf("results are %+v, len is %d", events, len(events))
	return events, err
}
