package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/migotom/cell-centre-services/pkg/entities"
	"github.com/migotom/cell-centre-services/pkg/pb"
)

type RoleRepositoryMock struct {
	mock.Mock
}

func (m *RoleRepositoryMock) Get(ctx context.Context, filter *pb.RoleFilter) (*entities.Role, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*entities.Role), args.Error(1)
}
