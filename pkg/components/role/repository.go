package role

import (
	"context"

	"github.com/migotom/cell-centre-backend/pkg/entities"
	pb "github.com/migotom/cell-centre-backend/pkg/pb"
)

type Repository interface {
	Get(tx context.Context, filter *pb.RoleFilter) (*entities.Role, error)
}
