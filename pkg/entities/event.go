package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventType string

const (
	NewEmployeeEvent    EventType = "NewEmployee"
	UpdateEmployeeEvent EventType = "UpdateEmployee"
	DeleteEmployeeEvent EventType = "DeleteEmployee"
)

// Event entity definition.
type Event struct {
	EventID       string `bson:"event_id"`
	Channel       string
	Type          EventType
	AggregateID   string
	AggregateType string
	Data          interface{}
	Originator    EventOriginator
	CreatedAt     time.Time
}

type EventOriginator struct {
	EntityID primitive.ObjectID `bson:"entity_id"`
	Entity   string
	Login    string
}
