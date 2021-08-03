// Package entities contains structs that reflect database entities
package entities

import "github.com/google/uuid"

// Cat contains all data related to cat and stored in database
type Cat struct {
	ID    uuid.UUID `bson:"_id"`
	Name  string    `bson:"name"`
	Color string    `bson:"color"`
	Age   int       `bson:"age"`
	Price float64   `bson:"price"`
}
