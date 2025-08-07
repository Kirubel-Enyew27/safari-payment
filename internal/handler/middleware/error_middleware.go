package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Kirubel-Enyew27/safari-payment/internal/errors"
	response2 "github.com/Kirubel-Enyew27/safari-payment/internal/model/resonse"

	"github.com/joomcode/errorx"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func ErrorHandler() gin.HandlerFunc {
	debugMode := strings.ToLower(os.Getenv("DEBUG"))
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			e := c.Errors[0]
			err := e.Unwrap()

			response := CastErrorResponse(err)
			if response != nil {
				er := errorx.Cast(err)
				if debugMode == "true" {
					response.Description = fmt.Sprintf("Error: %v", er)
					response.StackTrace = fmt.Sprintf("%+v", errorx.EnsureStackTrace(err))
				}
				response2.SendErrorResponse(c, response)
				return
			}

			response2.SendErrorResponse(c, &response2.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Unknown server error",
			})
			return
		}
	}
}
func ErrorFields(err error) []response2.FieldError {
	var errs []response2.FieldError

	if data, ok := err.(validation.Errors); ok {
		for i, v := range data {
			nestedErrors := ErrorFields(v)
			if len(nestedErrors) > 0 {
				errs = append(errs, nestedErrors...)
			} else {
				errs = append(errs, response2.FieldError{
					Name:        i,
					Description: v.Error(),
				})
			}
		}

		return errs
	}

	return nil
}

func CastErrorResponse(err error) *response2.ErrorResponse {
	for _, e := range errors.Error {
		if errorx.IsOfType(err, e.ErrorType) {
			er := errorx.Cast(err)
			response := response2.ErrorResponse{
				Code:       e.StatusCode,
				Message:    er.Message(),
				FieldError: ErrorFields(er.Cause()),
			}
			return &response
		}
	}
	return nil
}
