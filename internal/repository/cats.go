// Package repository encapsulates work with databases
package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/evleria/mongo-crud/internal/repository/entities"
)

var (
	// ErrNotFound means entity is not found in repository
	ErrNotFound = errors.New("not found")
)

// Cats contains methods for manipulating with cats collection
type Cats interface {
	Insert(ctx context.Context, name, color string, age int) (uuid.UUID, error)
	GetAll(ctx context.Context) ([]entities.Cat, error)
	GetOne(ctx context.Context, id uuid.UUID) (entities.Cat, error)
	Delete(ctx context.Context, id uuid.UUID) error
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

func (c *cats) GetOne(ctx context.Context, id uuid.UUID) (entities.Cat, error) {
	cat := entities.Cat{}
	err := c.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&cat)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return cat, ErrNotFound
	} else if err != nil {
		return cat, err
	}
	return cat, nil
}

func (c *cats) Delete(ctx context.Context, id uuid.UUID) error {
	if r, err := c.collection.DeleteOne(ctx, bson.M{"_id": id}); err != nil {
		return err
	} else if r.DeletedCount == 0 {
		return ErrNotFound
	}
	return nil
}