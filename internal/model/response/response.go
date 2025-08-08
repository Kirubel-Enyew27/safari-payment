package response

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Kirubel-Enyew27/safari-payment/internal/errors"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/joomcode/errorx"
)

type Response struct {
	// OK is only true if the request was successful.
	OK bool `json:"ok"`
	// MetaData contains additional data like filtering, pagination, etc.
	MetaData *MetaData `json:"meta_data,omitempty"`
	// Data contains the actual data of the response.
	Data interface{} `json:"data,omitempty"`
	// Error contains the error detail if the request was not successful.
	Error *ErrorResponse `json:"error,omitempty"`
}

type MetaData struct {
	// Total is the total number of data without pagination
	Total int `json:"total"`
	// Extra contains other response specific data
	Extra interface{} `json:"extra,omitempty"`
}

type ErrorResponse struct {
	// Code is the error code. It is not status code
	Code int `json:"code"`
	// Message is the error message.
	Message string `json:"message,omitempty"`
	// Description is the error description.
	Description string `json:"description,omitempty"`
	// StackTrace is the stack trace of the error.
	// It is only returned for debugging
	StackTrace string `json:"stack_trace,omitempty"`
	// FieldError is the error detail for each field, if available that is.
	FieldError []FieldError `json:"field_error,omitempty"`
}

type FieldError struct {
	// Name is the name of the field that caused the error.
	Name string `json:"name"`
	// Description is the error description for this field.
	Description string `json:"description"`
}

func SendSuccessResponse(ctx *gin.Context, statusCode int, data interface{},
	metaData *MetaData) {

	ctx.JSON(
		statusCode,
		Response{
			OK:       true,
			MetaData: metaData,
			Data:     data,
		},
	)
}

func SendErrorResponse(ctx *gin.Context, err *ErrorResponse) {
	ctx.AbortWithStatusJSON(err.Code, Response{
		OK:    false,
		Error: err,
	})
}

func GetErrorFrom(err error) *ErrorResponse {
	// debugMode := viper.GetBool("debug")
	debugMode := strings.ToLower(os.Getenv("DEBUG"))
	for _, e := range errors.Error {
		if errorx.IsOfType(err, e.ErrorType) {
			er := errorx.Cast(err)
			res := ErrorResponse{
				Code:       e.StatusCode,
				Message:    er.Message(),
				FieldError: ErrorFields(er.Cause()),
			}

			if debugMode == "true" {
				res.Description = fmt.Sprintf("Error: %v", er)
				res.StackTrace = fmt.Sprintf("%+v", errorx.EnsureStackTrace(err))
			}

			return &res
		}
	}

	return &ErrorResponse{
		Code:    http.StatusInternalServerError,
		Message: "Unknown server error",
	}
}

func ErrorFields(err error) []FieldError {
	var errs []FieldError

	if data, ok := err.(validation.Errors); ok {
		for i, v := range data {
			errs = append(errs, FieldError{
				Name:        i,
				Description: v.Error(),
			},
			)
		}

		return errs
	}

	return nil
}
