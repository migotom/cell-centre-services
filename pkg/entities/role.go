package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Role entity definition.
type Role struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name,omitempty"`
}
