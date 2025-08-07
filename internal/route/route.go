package route

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Router struct {
	Method      string
	Path        string
	Handler     gin.HandlerFunc
	Middlewares []gin.HandlerFunc
}

func RegisterRoute(
	grp *gin.RouterGroup,
	routes []Router,
	logger *zap.Logger,
) {
	for _, route := range routes {
		var handler []gin.HandlerFunc

		handler = append(handler, route.Middlewares...)
		handler = append(handler, route.Handler)

		grp.Handle(route.Method, route.Path, handler...)
	}
}
