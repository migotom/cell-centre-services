package entities

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TokenClaimer represents interface of entity that is able to provide auth claims parameters and login into system.
type TokenClaimer interface {
	GetID() primitive.ObjectID
	GetEntity() string
	GetLogin() string
	GetRoles() []Role
}

// TokenClaims repesents JWT authentication claims.
type TokenClaims struct {
	Entity             string             `json:"entity"`
	EntityID           primitive.ObjectID `json:"entity_id"`
	Login              string             `json:"login"`
	Roles              []string           `json:"roles"`
	jwt.StandardClaims `bson:"-"`
}

// NewTokenClaims returns JWT claims for specified entity.
func NewTokenClaims(expiration time.Duration, entity TokenClaimer) *TokenClaims {
	claims := TokenClaims{
		Entity:   entity.GetEntity(),
		EntityID: entity.GetID(),
		Login:    entity.GetLogin(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiration).Unix(),
		},
	}
	for _, role := range entity.GetRoles() {
		claims.Roles = append(claims.Roles, role.Name)
	}
	return &claims
}

// HasRole verifies that claims has at least one of given roles.
func (claims *TokenClaims) HasRole(authorizedRoles []string) bool {
	for _, authorizedRole := range authorizedRoles {
		for _, role := range claims.Roles {
			if role == authorizedRole {
				return true
			}
		}
	}
	return false
}
