package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/migotom/cell-centre-backend/pkg/entities"
	"github.com/migotom/cell-centre-backend/pkg/pb"
)

type EmployeRepositoryMock struct {
	mock.Mock
}

func (m *EmployeRepositoryMock) Get(ctx context.Context, filter *pb.EmployeeFilter) (*entities.Employee, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*entities.Employee), args.Error(1)
}
func (m *EmployeRepositoryMock) New(ctx context.Context, request *entities.Employee) (*entities.Employee, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(*entities.Employee), args.Error(1)
}
func (m *EmployeRepositoryMock) Update(ctx context.Context, request *entities.Employee) (*entities.Employee, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(*entities.Employee), args.Error(1)
}
func (m *EmployeRepositoryMock) Delete(ctx context.Context, filter *pb.EmployeeFilter) error {
	args := m.Called(ctx, filter)
	return args.Error(0)
}
