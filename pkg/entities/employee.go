package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const entityType = "employee"

// Employee entity definition.
type Employee struct {
	ID        primitive.ObjectID `bson:"_id"`
	Email     string             `bson:"email,omitempty"`
	Password  string             `bson:"password,omitempty" json:"-"`
	Name      string             `bson:"name,omitempty"`
	Phone     string             `bson:"phone,omitempty"`
	CreatedAt *time.Time         `bson:"created_at,omitempty"`
	UpdatedAt *time.Time         `bson:"updated_at,omitempty"`
	Roles     []Role             `bson:"roles,omitempty"`
}

// GetEntity returns type employee's type of entity.
func (employee *Employee) GetEntity() string {
	return entityType
}

// GetID returns employee's ID.
func (employee *Employee) GetID() primitive.ObjectID {
	return employee.ID
}

// GetLogin returns employee's email as Login.
func (employee *Employee) GetLogin() string {
	return employee.Email
}

// GetRoles returns employee's roles.
func (employee *Employee) GetRoles() []Role {
	return employee.Roles
}
