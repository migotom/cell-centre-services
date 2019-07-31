package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	"github.com/migotom/cell-centre-services/pkg/entities"
)

var (
	JwtSecret       = []byte("upersecretpass")
	TokenExpiration = 60 * time.Minute
)

// NewToken return new token for given employee.
func NewToken(entity entities.TokenClaimer) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, entities.NewTokenClaims(TokenExpiration, entity)).SignedString(JwtSecret)
}

// ParseToken parses given token and returns token claims with embedded JWT claims if token is valid.
func ParseToken(authenticate string) (entities.TokenClaims, error) {
	var claims entities.TokenClaims
	token, err := jwt.ParseWithClaims(authenticate, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, AuthError{Reason: ErrDecryptionToken}
		}
		return JwtSecret, nil
	})
	if err != nil {
		return entities.TokenClaims{}, AuthError{Reason: ErrDecryptionToken, Err: err}
	}

	if !token.Valid {
		return entities.TokenClaims{}, AuthError{Reason: ErrInvalidToken}
	}

	return claims, nil
}

// ValidPassword verifies given password with hashed one.
func ValidPassword(hashedPwd string, plainPwd string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd)); err != nil {
		return false
	}
	return true
}
