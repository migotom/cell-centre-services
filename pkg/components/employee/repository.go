package employee

import (
	"context"

	"github.com/migotom/cell-centre-services/pkg/entities"
	pb "github.com/migotom/cell-centre-services/pkg/pb"
)

// Repository of emploeyee.
type Repository interface {
	Get(ctx context.Context, filter *pb.EmployeeFilter) (*entities.Employee, error)
	New(ctx context.Context, request *entities.Employee) (*entities.Employee, error)
	Update(ctx context.Context, request *entities.Employee) (*entities.Employee, error)
	Delete(ctx context.Context, filter *pb.EmployeeFilter) error
}
