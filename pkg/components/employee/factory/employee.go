package factory

import (
	"context"

	"github.com/golang/protobuf/ptypes"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	"go.mongodb.org/mongo-driver/bson/primitive"

	roleRepository "github.com/migotom/cell-centre-backend/pkg/components/role"
	"github.com/migotom/cell-centre-backend/pkg/entities"
	"github.com/migotom/cell-centre-backend/pkg/pb"
)

// EmployeeEntityFactory is entities.Employee factory.
type EmployeeEntityFactory struct {
	role roleRepository.Repository
}

// NewEmployeeEntityFactory creates new factory.
func NewEmployeeEntityFactory(roleRepository roleRepository.Repository) *EmployeeEntityFactory {
	return &EmployeeEntityFactory{
		role: roleRepository,
	}
}

// NewFromNewEmployeeRequest creates Employee entity from NewEmployeeRequest message.
func (factory *EmployeeEntityFactory) NewFromNewEmployeeRequest(e *pb.NewEmployeeRequest) (employee *entities.Employee, err error) {
	if employee, err = factory.newEntityFromBase(e); err != nil {
		return
	}

	employee.ID = primitive.NewObjectID()
	return
}

// NewFromUpdateEmployeeRequest creates Employee entity from UpdateEmployeeRequest message.
func (factory *EmployeeEntityFactory) NewFromUpdateEmployeeRequest(e *pb.UpdateEmployeeRequest) (employee *entities.Employee, err error) {
	if employee, err = factory.newEntityFromBase(e); err != nil {
		return
	}

	if employee.ID, err = primitive.ObjectIDFromHex(e.Id); err != nil {
		return
	}

	return
}

// NewFromEmployeeFilter creates Employee entity from EmployeeFilter message.
func (factory *EmployeeEntityFactory) NewFromEmployeeFilter(e *pb.EmployeeFilter) (employee *entities.Employee, err error) {
	employee = &entities.Employee{}

	if e.GetId() != "" {
		employee.ID, _ = primitive.ObjectIDFromHex(e.Id)
	}
	if e.GetEmail() != "" {
		employee.Email = e.Email
	}

	return
}

// NewFromEmployee creates Employee entity from Employee message.
func (factory *EmployeeEntityFactory) NewFromEmployee(e *pb.Employee) (employee *entities.Employee, err error) {
	if employee, err = factory.newEntityFromBase(e); err != nil {
		return
	}

	if employee.ID, err = primitive.ObjectIDFromHex(e.Id); err != nil {
		return
	}

	if e.CreatedAt != nil {
		createdAt, err := ptypes.Timestamp(e.CreatedAt)
		if err != nil {
			return nil, err
		}
		employee.CreatedAt = &createdAt
	}

	if e.UpdatedAt != nil {
		updatedAt, err := ptypes.Timestamp(e.UpdatedAt)
		if err != nil {
			return nil, err
		}
		employee.UpdatedAt = &updatedAt
	}

	return
}

func (factory *EmployeeEntityFactory) newEntityFromBase(e baseEmployeeInterface) (*entities.Employee, error) {
	employee := entities.Employee{
		Email:    e.GetEmail(),
		Password: e.GetPassword(),
		Name:     e.GetName(),
		Phone:    e.GetPhone(),
	}

	for _, role := range e.GetRoles() {
		filter := pb.RoleFilter{
			Name: role.GetName(),
			Id:   role.GetId(),
		}

		if factory.role != nil {
			// verification of role existence is required
			entityRole, err := factory.role.Get(context.Background(), &filter)
			if err != nil {
				return nil, err
			}
			employee.Roles = append(employee.Roles, *entityRole)
		} else {
			// verification of role existence doesn't need
			entityRole := entities.Role{Name: filter.Name}
			var err error
			if entityRole.ID, err = primitive.ObjectIDFromHex(filter.Id); err != nil {
				return nil, err
			}
			employee.Roles = append(employee.Roles, entityRole)
		}
	}
	return &employee, nil
}

type employeeInterface interface {
	GetId() string
	GetCreatedAt() *timestamp.Timestamp
	GetUpdatedAt() *timestamp.Timestamp
	baseEmployeeInterface
}

type baseEmployeeInterface interface {
	GetEmail() string
	GetName() string
	GetPassword() string
	GetPhone() string
	GetRoles() []*pb.Role
}
