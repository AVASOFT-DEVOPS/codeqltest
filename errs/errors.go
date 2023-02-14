package errs

import (
	"fmt"
	"golang/dto"
	"net/http"
)

type AppError struct {
	StatusCode      int    `json:",omitempty"`
	Message         string `json:"message"`
	ErrorId         int    `json:"errorId"`
	Code            string
	ErrorValidation error
	Errors          []dto.ValidateResDTO
}

type AppError1 struct {
	StatusCode int    `json:"omitempty"`
	Message    string `json:"omitempty"`
	ErrorId    int    `json:"omitempty"`
	Code       string `json:"omitempty"`
	Errors     []dto.ValidateResponse
}
type AppErrorvalidation struct {
	Code    int                  `json:",omitempty"`
	Errors  []dto.ValidateResDTO `json:",omitempty"`
	Message string               `json:",omitempty"`
}

func (e AppError) Error() string {
	//TODO implement me
	panic("implement me")
}

func (e AppError) AsMessage() *AppError {
	return &AppError{
		Message: e.Message,
	}
}

func NewUnexpectedError(message string, code string) *AppError {
	return &AppError{
		Message:    message,
		Code:       code,
		StatusCode: http.StatusInternalServerError,
	}
}

func ValidateResponse(Error []dto.ValidateResDTO, code int, message string) *AppErrorvalidation {
	return &AppErrorvalidation{
		Errors:  Error,
		Code:    code,
		Message: message,
	}
}

type ErrorResponse struct {
	Errors struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
}

type ErrorRes struct {
	Errors struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
}
type ValidateResDTO struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

//createTL
type CUserError struct {
	Code            int                  `json:"Code,omitempty"`
	Message         string               `json:"Message"`
	ErrorId         string               `json:"ErrorId"`
	ValidationError []dto.ValidateResDTO `json:"ValidationError,omitempty"` //Status  string `json:"status"`
}

func UserValidation(Error string) *CUserError {
	return &CUserError{
		Message: Error,
		Code:    http.StatusBadRequest,
	}
}

func UserRequestValidation(Error []dto.ValidateResDTO) *CUserError {
	fmt.Println("UserRequestValidation Error", Error)
	fmt.Println("UserRequestValidation struct", &CUserError{
		Message:         "Request Validation Failed",
		Code:            http.StatusBadRequest,
		ErrorId:         "ADMREQ10002",
		ValidationError: Error,
	})
	return &CUserError{
		Message:         "Request Validation Failed",
		Code:            http.StatusBadRequest,
		ValidationError: Error,
	}
}
func UserAuthentication(message string) *CUserError {
	return &CUserError{
		ErrorId: "ADM001",
		Message: message,
		Code:    http.StatusForbidden,
	}
}

func UserNewUnexpectedError(message string) *CUserError {
	return &CUserError{
		Message: message,
		Code:    http.StatusInternalServerError,
	}
}

func UserNewNotFoundError(message string) *CUserError {
	return &CUserError{
		Message: message,
		Code:    http.StatusNotFound,
		ErrorId: "203",
		//Status:  "Failed",
	}
}

func (s CUserError) AsUserMessage() *CUserError{
	return &CUserError{
		ErrorId: "L001",
		Message: s.Message,
		ValidationError: s.ValidationError,
	}
}
func (e CUserError) Error() string {
	//TODO implement me
	panic("implement me")
}
