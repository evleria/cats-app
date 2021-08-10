// Package repository encapsulates work with databases
package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/evleria/cats-app/internal/repository/entities"
)

var (
	// ErrNotFound means entity is not found in repository
	ErrNotFound = errors.New("not found")
)

// Cats contains methods for manipulating with cats collection
type Cats interface {
	Insert(ctx context.Context, name, color string, age int, price float64) (uuid.UUID, error)
	GetAll(ctx context.Context) ([]entities.Cat, error)
	GetOne(ctx context.Context, id uuid.UUID) (entities.Cat, error)
	Delete(ctx context.Context, id uuid.UUID) error
	UpdatePrice(ctx context.Context, id uuid.UUID, price float64) error
}

type cats struct {
	collection *mongo.Collection
}

// NewCatsRepository creates new cats repository
func NewCatsRepository(mongoDB *mongo.Database) Cats {
	return &cats{
		collection: mongoDB.Collection("cats"),
	}
}

func (c *cats) Insert(ctx context.Context, name, color string, age int, price float64) (uuid.UUID, error) {
	cat := entities.Cat{
		ID:    uuid.New(),
		Name:  name,
		Color: color,
		Age:   age,
		Price: price,
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

func (c *cats) UpdatePrice(ctx context.Context, id uuid.UUID, price float64) error {
	if r, err := c.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"price": price}}); err != nil {
		return err
	} else if r.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}
