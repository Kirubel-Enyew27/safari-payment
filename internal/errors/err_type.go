package errors

import (
	"net/http"

	"github.com/joomcode/errorx"
)

type ErrorType struct {
	StatusCode int
	ErrorType  *errorx.Type
}

var Error = []ErrorType{
	{
		StatusCode: http.StatusInternalServerError,
		ErrorType:  ErrInternalServerError,
	},
	{
		StatusCode: http.StatusInternalServerError,
		ErrorType:  ErrUnableToGet,
	},
	{
		StatusCode: http.StatusInternalServerError,
		ErrorType:  ErrUnableTocreate,
	},

	{
		StatusCode: http.StatusBadRequest,
		ErrorType:  ErrDataAlredyExist,
	},
	{
		StatusCode: http.StatusNotFound,
		ErrorType:  ErrNoRecordFound,
	},
	{
		StatusCode: http.StatusInternalServerError,
		ErrorType:  ErrFailedPayment,
	},
	{
		StatusCode: http.StatusBadRequest,
		ErrorType:  ErrBadRequest,
	},
	{
		StatusCode: http.StatusInternalServerError,
		ErrorType:  UnexpectedError,
	},
	{
		StatusCode: http.StatusInternalServerError,
		ErrorType:  ErrCreateRequest,
	},
	{
		StatusCode: http.StatusUnauthorized,
		ErrorType:  ErrAccessToken,
	},
}

// list of error namespaces
var (
	databaseError    = errorx.NewNamespace("database error")
	resourceNotFound = errorx.NewNamespace("not found")
	unauthorized     = errorx.NewNamespace("unable to get access token")
	AccessDenied     = errorx.RegisterTrait("You are not authorized to perform the action")
	Ineligible       = errorx.RegisterTrait("You are not eligible to perform the action")
	serverError      = errorx.NewNamespace("server error")
	badRequest       = errorx.NewNamespace("bad request error")
	requestError     = errorx.NewNamespace("error while communicating with safari")
	paymentErr       = errorx.NewNamespace("error while accepting payment")
)

// list of errors types in all of the above namespaces

var (
	ErrUnableTocreate      = errorx.NewType(databaseError, "unable to create")
	ErrDataAlredyExist     = errorx.NewType(databaseError, "data alredy exist")
	ErrUnableToGet         = errorx.NewType(databaseError, "unable to get")
	ErrInternalServerError = errorx.NewType(serverError, "internal server error")
	ErrUnExpectedError     = errorx.NewType(serverError, "unexpected error occurred")
	ErrNoRecordFound       = errorx.NewType(resourceNotFound, "no record found")
	ErrBadRequest          = errorx.NewType(badRequest, "bad request error")
	UnexpectedError        = errorx.NewType(serverError, "invalid value")
	ErrCreateRequest       = errorx.NewType(requestError, "failed to create request")
	ErrFailedPayment       = errorx.NewType(paymentErr, "payment faled")
	ErrAccessToken         = errorx.NewType(unauthorized, "unauthorized")
)
