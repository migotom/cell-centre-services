package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/migotom/cell-centre-services/pkg/components/role"
	"github.com/migotom/cell-centre-services/pkg/entities"
	"github.com/migotom/cell-centre-services/pkg/pb"
)

type roleRepo struct {
	DB *mongo.Database
}

func NewRoleRepository(db *mongo.Database) role.Repository {
	return &roleRepo{
		DB: db,
	}
}

func (repository *roleRepo) fetchOne(ctx context.Context, filter bson.D) (*entities.Role, error) {
	collection := repository.DB.Collection("roles")
	res := collection.FindOne(ctx, filter)

	var role entities.Role
	if err := res.Decode(&role); err != nil {
		return nil, err
	}
	return &role, nil
}

// Get return employee entity by ID.
func (repository *roleRepo) Get(ctx context.Context, filter *pb.RoleFilter) (*entities.Role, error) {
	switch {
	case filter.GetId() != "":
		id, err := primitive.ObjectIDFromHex(filter.GetId())
		if err != nil {
			return nil, err
		}
		return repository.fetchOne(ctx, bson.D{{"_id", id}})
	case filter.GetName() != "":
		return repository.fetchOne(ctx, bson.D{{"name", filter.GetName()}})
	}
	return nil, errors.New("unknown role filter")
}
