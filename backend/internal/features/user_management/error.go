package usermanagement

import "fmt"

type UserManagementDomainError error

var (
	ErrWrongCredentials          = fmt.Errorf("wrongCredentials")
	ErrUserNotFound              = fmt.Errorf("userNotFound")
	ErrEmailAlreadyExists        = fmt.Errorf("emailAlreadyExists")
	ErrInvalidTokenSigningMethod = fmt.Errorf("")
	ErrInvalidClaimsFormat       = fmt.Errorf("")
	ErrTokenExpired              = fmt.Errorf("")
)
