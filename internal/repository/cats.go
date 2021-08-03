// Package repository encapsulates work with databases
package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/evleria/mongo-crud/internal/repository/entities"
)

// Cats contains methods for manipulating with cats collection
type Cats interface {
	Insert(ctx context.Context, name, color string, age int) (uuid.UUID, error)
	GetAll(ctx context.Context) ([]entities.Cat, error)
}

type cats struct {
	collection *mongo.Collection
}

// NewCatsRepository creates new cats repository
func NewCatsRepository(mongoClient *mongo.Client, dbName string) Cats {
	return &cats{
		collection: mongoClient.Database(dbName).Collection("cats"),
	}
}

func (c *cats) Insert(ctx context.Context, name, color string, age int) (uuid.UUID, error) {
	cat := entities.Cat{
		ID:    uuid.New(),
		Name:  name,
		Color: color,
		Age:   age,
	}

	_, err := c.collection.InsertOne(ctx, cat)
	if err != nil {
		return uuid.UUID{}, err
	}
	return cat.ID, err
}

func (c *cats) GetAll(ctx context.Context) ([]entities.Cat, error) {
	cursor, err := c.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	result := []entities.Cat{}

	for cursor.Next(ctx) {
		cat := new(entities.Cat)
		if err := cursor.Decode(cat); err != nil {
			return nil, err
		}
		result = append(result, *cat)
	}

	if err := cursor.Close(ctx); err != nil {
		return nil, err
	}
	return result, nil
}
