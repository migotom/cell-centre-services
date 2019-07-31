package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"google.golang.org/grpc/metadata"

	"github.com/migotom/cell-centre-services/pkg/components/auth"
	"github.com/migotom/cell-centre-services/pkg/entities"
	"github.com/migotom/cell-centre-services/pkg/helpers"
	"github.com/migotom/cell-centre-services/pkg/helpers/mocks"
	"github.com/migotom/cell-centre-services/pkg/pb"
)

func TestDefaultInterceptor(t *testing.T) {
	cases := []struct {
		Name           string
		TokenFunc      func() string
		ExpectedErr    string
		ExpectedClaims func(ctx context.Context)
	}{
		{
			Name: "Valid token",
			TokenFunc: func() string {
				token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, entities.TokenClaims{Roles: []string{"admin"}}).SignedString(auth.JwtSecret)
				return token
			},
			ExpectedErr: "",
			ExpectedClaims: func(ctx context.Context) {
				claims, ok := ctx.Value(ContextKeyClaims).(entities.TokenClaims)
				if !ok {
					t.Errorf("missing claims")
				}
				if claims.Roles[0] != "admin" {
					t.Errorf("expected to have admin role in claims")
				}
			},
		},
		{
			Name: "Invalid token",
			TokenFunc: func() string {
				return "xxxxxx"
			},
			ExpectedErr: "rpc error: code = Unauthenticated desc = Request unauthenticated with error: Error during token decryption (token contains an invalid number of segments)",
			ExpectedClaims: func(ctx context.Context) {
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			log, _ := zap.NewProduction()
			defer log.Sync()

			token := tc.TokenFunc()
			md := metadata.New(map[string]string{headerAuthorize: token})
			ctx := metadata.NewIncomingContext(context.Background(), md)
			delivery := NewAuthenticateDelivery(log, nil)

			newCtx, err := delivery.DefaultInterceptor(ctx)
			helpers.AssertErrors(t, err, tc.ExpectedErr)

			if newCtx != nil {
				tc.ExpectedClaims(newCtx)
			}
		})
	}
}

func TestAuthenticate(t *testing.T) {
	cases := []struct {
		Name                string
		AuthRequest         pb.AuthRequest
		ExpectedMockCalls   func(*mocks.EmployeRepositoryMock)
		ExpectedErr         string
		ExpectedClaimsRoles []string
	}{
		{
			Name: "Valid request",
			AuthRequest: pb.AuthRequest{
				Entity:   pb.AuthRequest_EMPLOYEE,
				Login:    "admin@page.com",
				Password: "test123",
			},
			ExpectedMockCalls: func(m *mocks.EmployeRepositoryMock) {
				m.On("Get", mock.Anything, &pb.EmployeeFilter{Email: "admin@page.com"}).
					Return(&entities.Employee{
						Email:    "admin@page.com",
						Password: helpers.HashPassword("test123"),
						Roles:    []entities.Role{{Name: "admin"}},
					}, nil)
			},
			ExpectedErr:         "",
			ExpectedClaimsRoles: []string{"admin"},
		},
		{
			Name: "Invalid credentials",
			AuthRequest: pb.AuthRequest{
				Entity:   pb.AuthRequest_EMPLOYEE,
				Login:    "nobody@page.com",
				Password: "test123",
			},
			ExpectedMockCalls: func(m *mocks.EmployeRepositoryMock) {
				m.On("Get", mock.Anything, &pb.EmployeeFilter{Email: "nobody@page.com"}).
					Return(&entities.Employee{}, errors.New("not existing"))
			},
			ExpectedErr:         "rpc error: code = Unauthenticated desc = Can't authenticate: Invalid credentials",
			ExpectedClaimsRoles: []string(nil),
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			log, _ := zap.NewProduction()
			defer log.Sync()

			employeeRepositoryMock := mocks.EmployeRepositoryMock{}
			tc.ExpectedMockCalls(&employeeRepositoryMock)

			delivery := NewAuthenticateDelivery(log, &employeeRepositoryMock)
			res, err := delivery.Authenticate(context.Background(), &tc.AuthRequest)
			helpers.AssertErrors(t, err, tc.ExpectedErr)
			if err != nil {
				return
			}

			employeeRepositoryMock.AssertExpectations(t)

			claims, err := auth.ParseToken(res.Token)
			helpers.AssertErrors(t, err, tc.ExpectedErr)
			if err != nil {
				return
			}
			assert.Equal(t, claims.Roles, tc.ExpectedClaimsRoles)
		})
	}
}
