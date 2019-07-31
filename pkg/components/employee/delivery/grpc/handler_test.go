package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"github.com/migotom/cell-centre-services/pkg/entities"
	"github.com/migotom/cell-centre-services/pkg/helpers"
	"github.com/migotom/cell-centre-services/pkg/helpers/mocks"
	"github.com/migotom/cell-centre-services/pkg/pb"
)

func TestGetEmployee(t *testing.T) {
	cases := []struct {
		Name              string
		Filter            pb.EmployeeFilter
		ExpectedMockCalls func(*mocks.EmployeRepositoryMock, *mocks.RoleRepositoryMock)
		ExpectedEmployee  *pb.Employee
		ExpectedErr       string
	}{
		{
			Name: "Valid request by ID",
			Filter: pb.EmployeeFilter{
				Id: "1",
			},
			ExpectedMockCalls: func(e *mocks.EmployeRepositoryMock, r *mocks.RoleRepositoryMock) {
				id, _ := primitive.ObjectIDFromHex("5d3783ee28ae9468bc528906")
				e.On("Get", mock.Anything, &pb.EmployeeFilter{Id: "1"}).
					Return(&entities.Employee{
						ID:       id,
						Email:    "admin@page.com",
						Password: "$2a$04$6UsCk8fCtstbKTT1fmgPa.SxO4L8BIxrjStRuXMiNTM9HdzIDdGBK",
						Roles:    []entities.Role{{ID: id, Name: "admin"}},
					}, nil)
			},
			ExpectedEmployee: &pb.Employee{
				Id:       "5d3783ee28ae9468bc528906",
				Email:    "admin@page.com",
				Password: "$2a$04$6UsCk8fCtstbKTT1fmgPa.SxO4L8BIxrjStRuXMiNTM9HdzIDdGBK",
				Roles:    []*pb.Role{&pb.Role{Id: "5d3783ee28ae9468bc528906", Name: "admin"}},
			},
		},
		{
			Name: "Valid request by Email",
			Filter: pb.EmployeeFilter{
				Email: "admin@page.com",
			},
			ExpectedMockCalls: func(e *mocks.EmployeRepositoryMock, r *mocks.RoleRepositoryMock) {
				id, _ := primitive.ObjectIDFromHex("5d3783ee28ae9468bc528906")
				e.On("Get", mock.Anything, &pb.EmployeeFilter{Email: "admin@page.com"}).
					Return(&entities.Employee{
						ID:       id,
						Email:    "admin@page.com",
						Password: "$2a$04$6UsCk8fCtstbKTT1fmgPa.SxO4L8BIxrjStRuXMiNTM9HdzIDdGBK",
						Roles:    []entities.Role{{ID: id, Name: "admin"}},
					}, nil)
			},
			ExpectedEmployee: &pb.Employee{
				Id:       "5d3783ee28ae9468bc528906",
				Email:    "admin@page.com",
				Password: "$2a$04$6UsCk8fCtstbKTT1fmgPa.SxO4L8BIxrjStRuXMiNTM9HdzIDdGBK",
				Roles:    []*pb.Role{&pb.Role{Id: "5d3783ee28ae9468bc528906", Name: "admin"}},
			},
		},
		{
			Name: "Invalid request (not found)",
			Filter: pb.EmployeeFilter{
				Email: "nobody@page.com",
			},
			ExpectedMockCalls: func(e *mocks.EmployeRepositoryMock, r *mocks.RoleRepositoryMock) {
				e.On("Get", mock.Anything, &pb.EmployeeFilter{Email: "nobody@page.com"}).
					Return(&entities.Employee{}, errors.New("not found"))
			},
			ExpectedEmployee: nil,
			ExpectedErr:      "rpc error: code = NotFound desc = Can't get employee: not found",
		},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			log, _ := zap.NewProduction()
			defer log.Sync()

			employeeRepositoryMock := mocks.EmployeRepositoryMock{}
			roleRepositoryMock := mocks.RoleRepositoryMock{}

			tc.ExpectedMockCalls(&employeeRepositoryMock, &roleRepositoryMock)

			delivery := NewEmployeeDelivery(
				log,
				&employeeRepositoryMock,
				&roleRepositoryMock,
				nil,
			)
			employee, err := delivery.GetEmployee(context.Background(), &tc.Filter)
			helpers.AssertErrors(t, tc.ExpectedErr, err)

			assert.Equal(t, tc.ExpectedEmployee, employee)
		})
	}
}

func TestNewEmployee(t *testing.T) {
	cases := []struct {
		Name              string
		Request           pb.NewEmployeeRequest
		ExpectedMockCalls func(*mocks.EmployeRepositoryMock, *mocks.RoleRepositoryMock)
		ExpectedEmployee  *pb.Employee
		ExpectedErr       string
	}{
		{
			Name: "Valid request",
			Request: pb.NewEmployeeRequest{
				Email:    "admin@page.com",
				Password: "$2a$04$6UsCk8fCtstbKTT1fmgPa.SxO4L8BIxrjStRuXMiNTM9HdzIDdGBK",
				Roles:    []*pb.Role{&pb.Role{Name: "admin"}},
			},
			ExpectedMockCalls: func(e *mocks.EmployeRepositoryMock, r *mocks.RoleRepositoryMock) {
				id, _ := primitive.ObjectIDFromHex("5d3783ee28ae9468bc528906")
				r.On("Get", mock.Anything, &pb.RoleFilter{Name: "admin"}).Return(&entities.Role{Name: "admin"}, nil)
				e.On("New", mock.Anything, mock.Anything).Return(&entities.Employee{
					ID:    id,
					Email: "admin@page.com",
				}, nil)
			},
			ExpectedEmployee: &pb.Employee{
				Id:    "5d3783ee28ae9468bc528906",
				Email: "admin@page.com",
			},
		},
		{
			Name: "Invalid request - missing role",
			Request: pb.NewEmployeeRequest{
				Email:    "admin@page.com",
				Password: "$2a$04$6UsCk8fCtstbKTT1fmgPa.SxO4L8BIxrjStRuXMiNTM9HdzIDdGBK",
			},
			ExpectedMockCalls: func(e *mocks.EmployeRepositoryMock, r *mocks.RoleRepositoryMock) {
			},
			ExpectedEmployee: &pb.Employee{},
			ExpectedErr:      "rpc error: code = InvalidArgument desc = Invalid request: Invalid employee roles",
		},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			log, _ := zap.NewProduction()
			defer log.Sync()

			employeeRepositoryMock := mocks.EmployeRepositoryMock{}
			roleRepositoryMock := mocks.RoleRepositoryMock{}

			tc.ExpectedMockCalls(&employeeRepositoryMock, &roleRepositoryMock)

			delivery := NewEmployeeDelivery(
				log,
				&employeeRepositoryMock,
				&roleRepositoryMock,
				nil,
			)
			employee, err := delivery.NewEmployee(context.Background(), &tc.Request)
			helpers.AssertErrors(t, tc.ExpectedErr, err)

			assert.Equal(t, tc.ExpectedEmployee, employee)
		})
	}
}

func TestUpdateEmployee(t *testing.T) {
	cases := []struct {
		Name              string
		Request           pb.UpdateEmployeeRequest
		ExpectedMockCalls func(*mocks.EmployeRepositoryMock, *mocks.RoleRepositoryMock)
		ExpectedEmployee  *pb.Employee
		ExpectedErr       string
	}{
		{
			Name: "Valid request",
			Request: pb.UpdateEmployeeRequest{
				Id:    "5d3783ee28ae9468bc528906",
				Email: "admin@page.com",
			},
			ExpectedMockCalls: func(e *mocks.EmployeRepositoryMock, r *mocks.RoleRepositoryMock) {
				id, _ := primitive.ObjectIDFromHex("5d3783ee28ae9468bc528906")
				r.On("Get", mock.Anything, &pb.RoleFilter{Name: "admin"}).Return(&entities.Role{Name: "admin"}, nil)
				e.On("Update", mock.Anything, &entities.Employee{
					ID:    id,
					Email: "admin@page.com",
				}).Return(&entities.Employee{
					ID:    id,
					Email: "admin@page.com",
				}, nil)
			},
			ExpectedEmployee: &pb.Employee{
				Id:    "5d3783ee28ae9468bc528906",
				Email: "admin@page.com",
			},
		},
		{
			Name: "Invalid request",
			Request: pb.UpdateEmployeeRequest{
				Id:    "5d3783ee28ae9468bc520101",
				Email: "admin@page.com",
			},
			ExpectedMockCalls: func(e *mocks.EmployeRepositoryMock, r *mocks.RoleRepositoryMock) {
				id, _ := primitive.ObjectIDFromHex("5d3783ee28ae9468bc520101")
				r.On("Get", mock.Anything, &pb.RoleFilter{Name: "admin"}).Return(&entities.Role{Name: "admin"}, nil)
				e.On("Update", mock.Anything, &entities.Employee{
					ID:    id,
					Email: "admin@page.com",
				}).Return(&entities.Employee{}, errors.New("missing"))
			},
			ExpectedEmployee: nil,
			ExpectedErr:      "rpc error: code = InvalidArgument desc = Can't update employee: Invalid employee data (missing)",
		},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			log, _ := zap.NewProduction()
			defer log.Sync()

			employeeRepositoryMock := mocks.EmployeRepositoryMock{}
			roleRepositoryMock := mocks.RoleRepositoryMock{}

			tc.ExpectedMockCalls(&employeeRepositoryMock, &roleRepositoryMock)

			delivery := NewEmployeeDelivery(
				log,
				&employeeRepositoryMock,
				&roleRepositoryMock,
				nil,
			)
			employee, err := delivery.UpdateEmployee(context.Background(), &tc.Request)
			helpers.AssertErrors(t, tc.ExpectedErr, err)

			assert.Equal(t, tc.ExpectedEmployee, employee)
		})
	}
}

func TestDeleteEmployee(t *testing.T) {
	cases := []struct {
		Name              string
		Filter            pb.EmployeeFilter
		ExpectedMockCalls func(*mocks.EmployeRepositoryMock, *mocks.RoleRepositoryMock)
		ExpectedErr       string
	}{
		{
			Name: "Valid request by ID",
			Filter: pb.EmployeeFilter{
				Id: "1",
			},
			ExpectedMockCalls: func(e *mocks.EmployeRepositoryMock, r *mocks.RoleRepositoryMock) {
				e.On("Delete", mock.Anything, &pb.EmployeeFilter{Id: "1"}).Return(nil)
			},
		},
		{
			Name: "Invalid request (not found)",
			Filter: pb.EmployeeFilter{
				Email: "nobody@page.com",
			},
			ExpectedMockCalls: func(e *mocks.EmployeRepositoryMock, r *mocks.RoleRepositoryMock) {
				e.On("Delete", mock.Anything, &pb.EmployeeFilter{Email: "nobody@page.com"}).
					Return(errors.New("not found"))
			},
			ExpectedErr: "rpc error: code = NotFound desc = Can't delete employee: Invalid employee data (not found)",
		},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			log, _ := zap.NewProduction()
			defer log.Sync()

			employeeRepositoryMock := mocks.EmployeRepositoryMock{}
			roleRepositoryMock := mocks.RoleRepositoryMock{}

			tc.ExpectedMockCalls(&employeeRepositoryMock, &roleRepositoryMock)

			delivery := NewEmployeeDelivery(
				log,
				&employeeRepositoryMock,
				&roleRepositoryMock,
				nil,
			)
			_, err := delivery.DeleteEmployee(context.Background(), &tc.Filter)
			helpers.AssertErrors(t, tc.ExpectedErr, err)
		})
	}
}
