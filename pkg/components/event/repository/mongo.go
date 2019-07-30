package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/migotom/cell-centre-backend/pkg/components/event"
	"github.com/migotom/cell-centre-backend/pkg/entities"
)

type mongoEventRepo struct {
	DB *mongo.Database
}

// NewMongoEventRepository return new event MongoDB repository.
func NewMongoEventRepository(db *mongo.Database) event.Repository {
	return &mongoEventRepo{
		DB: db,
	}
}

func (repository *mongoEventRepo) New(ctx context.Context, event *entities.Event) error {
	collection := repository.DB.Collection("events")

	_, err := collection.InsertOne(ctx, event)
	if err != nil {
		return err
	}

	return nil
}
