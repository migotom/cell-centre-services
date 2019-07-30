package repository

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/migotom/cell-centre-backend/pkg/components/employee"
	"github.com/migotom/cell-centre-backend/pkg/entities"
	"github.com/migotom/cell-centre-backend/pkg/pb"
)

const collectionName = "employees"

type employeeRepository struct {
	DB *mongo.Database
}

// NewEmployeeRepository return new employee MongoDB repository.
func NewEmployeeRepository(db *mongo.Database) employee.Repository {
	return &employeeRepository{
		DB: db,
	}
}

// TODO add comments
func (repository *employeeRepository) New(ctx context.Context, request *entities.Employee) (*entities.Employee, error) {
	collection := repository.DB.Collection(collectionName)

	now := time.Now()
	request.CreatedAt = &now
	request.UpdatedAt = &now

	res, err := collection.InsertOne(ctx, request)
	if err != nil {
		return nil, err
	}

	return repository.fetchOne(ctx, bson.D{{"_id", res.InsertedID.(primitive.ObjectID)}})
}

// Get return employee entity by ID.
func (repository *employeeRepository) Get(ctx context.Context, filter *pb.EmployeeFilter) (*entities.Employee, error) {
	switch {
	case filter.GetId() != "":
		id, err := primitive.ObjectIDFromHex(filter.GetId())
		if err != nil {
			return nil, err
		}
		return repository.fetchOne(ctx, bson.D{{"_id", id}})
	case filter.GetEmail() != "":
		return repository.fetchOne(ctx, bson.D{{"email", filter.GetEmail()}})
	}
	return nil, errors.New("Unknown employee filter")
}

func (repository *employeeRepository) Update(ctx context.Context, request *entities.Employee) (*entities.Employee, error) {
	collection := repository.DB.Collection(collectionName)

	now := time.Now()
	request.UpdatedAt = &now

	_, err := collection.UpdateOne(ctx, bson.M{"_id": request.ID}, bson.M{"$set": request})
	if err != nil {
		return nil, err
	}

	return repository.fetchOne(ctx, bson.D{{"_id", request.ID}})
}

func (repository *employeeRepository) Delete(ctx context.Context, filter *pb.EmployeeFilter) error {
	switch {
	case filter.GetId() != "":
		id, err := primitive.ObjectIDFromHex(filter.GetId())
		if err != nil {
			return err
		}
		return repository.deleteOne(ctx, bson.D{{"_id", id}})
	case filter.GetEmail() != "":
		return repository.deleteOne(ctx, bson.D{{"email", filter.GetEmail()}})
	}
	return errors.New("Unknown employee filter")
}

func (repository *employeeRepository) fetchOne(ctx context.Context, filter bson.D) (*entities.Employee, error) {
	collection := repository.DB.Collection(collectionName)
	res := collection.FindOne(ctx, filter)

	var employee entities.Employee
	if err := res.Decode(&employee); err != nil {
		return nil, err
	}

	return &employee, nil
}

func (repository *employeeRepository) deleteOne(ctx context.Context, filter bson.D) error {
	collection := repository.DB.Collection(collectionName)
	_, err := collection.DeleteOne(ctx, filter)

	return err
}
