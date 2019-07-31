package auth

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"

	"github.com/migotom/cell-centre-services/pkg/entities"
	"github.com/migotom/cell-centre-services/pkg/helpers"
)

func TestNewTokenClaims(t *testing.T) {
	cases := []struct {
		Name           string
		Claimer        entities.TokenClaimer
		ExpectedClaims *entities.TokenClaims
	}{
		{
			Name: "Employee with admin role",
			Claimer: &entities.Employee{
				Email: "admin@page.com",
				Roles: []entities.Role{{Name: "admin"}},
			},
			ExpectedClaims: &entities.TokenClaims{
				Entity: "employee",
				Login:  "admin@page.com",
				Roles:  []string{"admin"},
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(TokenExpiration).Unix(),
				},
			},
		},
		{
			Name: "Employee with no role",
			Claimer: &entities.Employee{
				Email: "nobody@page.com",
				Roles: nil,
			},
			ExpectedClaims: &entities.TokenClaims{
				Entity: "employee",
				Login:  "nobody@page.com",
				Roles:  []string(nil),
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(TokenExpiration).Unix(),
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			claims := entities.NewTokenClaims(TokenExpiration, tc.Claimer)
			assert.Equal(t, tc.ExpectedClaims, claims)
		})
	}
}

func TestParseToken(t *testing.T) {
	cases := []struct {
		Name           string
		AuthToken      func() string
		ExpectedErr    string
		ExpectedClaims entities.TokenClaims
	}{
		{
			Name: "Valid token",
			AuthToken: func() string {
				token, _ := NewToken(&entities.Employee{
					Email: "admin@page.com",
					Roles: []entities.Role{{Name: "admin"}},
				})
				return token
			},
			ExpectedErr: "",
			ExpectedClaims: entities.TokenClaims{
				Entity: "employee",
				Login:  "admin@page.com",
				Roles:  []string{"admin"},
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(TokenExpiration).Unix(),
				},
			},
		},
		{

			Name: "Expired token",
			AuthToken: func() string {
				token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, entities.TokenClaims{
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(TokenExpiration * -1).Unix(),
					},
				}).SignedString(JwtSecret)
				return token
			},
			ExpectedErr:    "Error during token decryption (token is expired by 1h0m0s)",
			ExpectedClaims: entities.TokenClaims{},
		},
		{
			Name: "Broken token string",
			AuthToken: func() string {
				return "xxdoajoiasjdoiajioasj"
			},
			ExpectedErr:    "Error during token decryption (token contains an invalid number of segments)",
			ExpectedClaims: entities.TokenClaims{},
		},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			token := tc.AuthToken()
			claims, err := ParseToken(token)
			helpers.AssertErrors(t, err, tc.ExpectedErr)
			assert.Equal(t, tc.ExpectedClaims, claims)
		})

	}
}
