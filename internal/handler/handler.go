package handler

import "github.com/gin-gonic/gin"

type Payment interface {
	AcceptPayment(c *gin.Context)
}
