package factory

import (
	"github.com/golang/protobuf/ptypes"

	"github.com/migotom/cell-centre-backend/pkg/entities"
	"github.com/migotom/cell-centre-backend/pkg/pb"
)

// EmployeePbFactory is pb.Employee factory.
type EmployeePbFactory struct {
}

// NewEmployeePbFactory creates new factory.
func NewEmployeePbFactory() *EmployeePbFactory {
	return &EmployeePbFactory{}
}

// NewFromEmployee creates new pb.Employee instance from Employee entity.
func (factory *EmployeePbFactory) NewFromEmployee(e *entities.Employee) (*pb.Employee, error) {
	employee := pb.Employee{
		Id:       e.ID.Hex(),
		Email:    e.Email,
		Password: e.Password,
		Name:     e.Name,
		Phone:    e.Phone,
	}

	createdAt, err := ptypes.TimestampProto(*e.CreatedAt)
	if err != nil {
		return nil, err
	}
	if createdAt != nil {
		employee.CreatedAt = createdAt
	}

	updatedAt, err := ptypes.TimestampProto(*e.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if updatedAt != nil {
		employee.UpdatedAt = updatedAt
	}

	for _, role := range e.Roles {
		employee.Roles = append(employee.Roles, &pb.Role{
			Id:   role.ID.Hex(),
			Name: role.Name,
		})
	}
	return &employee, nil
}
