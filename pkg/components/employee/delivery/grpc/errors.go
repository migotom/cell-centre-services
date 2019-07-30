package grpc

import "fmt"

type EmployeeDeliveryErrorReason int

const (
	ErrInvalidEmployeeData = iota
	ErrInvalidEmployeeRoles
	ErrInternal
)

type EmployeeDeliveryError struct {
	Reason EmployeeDeliveryErrorReason
	Err    error
}

func (err EmployeeDeliveryError) Error() string {
	if err.Err != nil {
		return fmt.Sprintf("%s (%v)", err.description(), err.Err)
	}
	return err.description()
}

func (err EmployeeDeliveryError) description() string {
	switch err.Reason {
	case ErrInvalidEmployeeData:
		return "Invalid employee data"

	case ErrInvalidEmployeeRoles:
		return "Invalid employee roles"

	case ErrInternal:
		return "Internal error"

	default:
		return "Unknown error"
	}
}
