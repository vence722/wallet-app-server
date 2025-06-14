package service

import "fmt"

type ServiceError struct {
	ErrType    string
	ErrMessage string
	Cause      error
}

func (se ServiceError) Error() string {
	return fmt.Sprintf("%s, message: %s, detail: %v", se.ErrType, se.ErrMessage, se.Cause)
}

func newServiceError(errType string, errMessage string, cause error) ServiceError {
	return ServiceError{ErrType: errType, ErrMessage: errMessage, Cause: cause}
}

const (
	ErrTypeRecordNotFound      = "record not found"
	ErrTypeInvalidRequestBody  = "invalid request body"
	ErrTypeInternalServerError = "internal server error"
)

const (
	ErrMessageUserNotFound         = "user not found"
	ErrMessageDBError              = "database error"
	ErrMessagePasswordNotValid     = "password is not valid"
	ErrMessageNegativeOrZeroAmount = "amount must be positive"
	ErrMessageInsufficientBalance  = "insufficient balance"
)
