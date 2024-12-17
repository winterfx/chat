package chat

import "fmt"

type ChatError struct {
	Code    int
	Message string
	Err     error
}

func (e *ChatError) Unwrap() error {
	return e.Err
}

func (e *ChatError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("Error %d: %s - %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

func (e *ChatError) WithError(err error) *ChatError {
	e.Err = err
	return e
}

var (
	ErrDBOperation = &ChatError{Code: 1004, Message: "db operation failed"}
	ErrPubMessage  = &ChatError{Code: 1005, Message: "publish message failed"}
	ErrSubMessage  = &ChatError{Code: 1006, Message: "subscribe message failed"}
)

func NewChatError(code int, message string, err error) *ChatError {
	return &ChatError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
