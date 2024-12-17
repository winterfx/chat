package user

import "fmt"

type UsersError struct {
	Code    int
	Message string
	Err     error
}

func (e *UsersError) Unwrap() error {
	return e.Err
}

func (e *UsersError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("Error %d: %s - %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

func (e *UsersError) WithError(err error) *UsersError {
	e.Err = err
	return e
}

var (
	ErrUserAlreadyRegistered = &UsersError{Code: 1001, Message: "user already registered"}
	ErrUserPasswordIncorrect = &UsersError{Code: 1002, Message: "incorrect password"}
	ErrUserNotFound          = &UsersError{Code: 1003, Message: "user not found"}
	ErrDBOperation           = &UsersError{Code: 1004, Message: "db operation failed"}

	ErrInvalidToken        = &UsersError{Code: 1005, Message: "invalid token"}
	ErrInvalidAuthProvider = &UsersError{Code: 1006, Message: "invalid auth provider"}
)

func NewUsersError(code int, message string, err error) *UsersError {
	return &UsersError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
