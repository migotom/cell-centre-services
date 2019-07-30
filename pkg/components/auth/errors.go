package auth

import (
	"fmt"
)

type AuthErrorReason int

const (
	ErrInvalidParameters = iota
	ErrInvalidCredentials
	ErrMissingToken
	ErrInvalidToken
	ErrDecryptionToken
	ErrInsufficientRights
)

type AuthError struct {
	Reason AuthErrorReason
	Err    error
}

func (err AuthError) Error() string {
	if err.Err != nil {
		return fmt.Sprintf("%s (%v)", err.description(), err.Err)
	}
	return err.description()
}

func (err AuthError) description() string {
	switch err.Reason {
	case ErrInvalidParameters:
		return "Invalid credentials parameters"

	case ErrInvalidCredentials:
		return "Invalid credentials"

	case ErrMissingToken:
		return "Missing token"

	case ErrInvalidToken:
		return "Invalid token"

	case ErrDecryptionToken:
		return "Error during token decryption"

	case ErrInsufficientRights:
		return "Insufficient rights"
	}
	return "Unknown error"
}
