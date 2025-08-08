package payment

import (
	"net/http"

	"github.com/Kirubel-Enyew27/safari-payment/internal/handler"
	"github.com/Kirubel-Enyew27/safari-payment/internal/handler/middleware"
	"github.com/Kirubel-Enyew27/safari-payment/internal/route"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitRoute(group *gin.RouterGroup, h handler.Payment,
	log *zap.Logger,
) {
	paymentRoutes := []route.Router{
		{
			Method:  http.MethodPost,
			Path:    "/payment/accept",
			Handler: h.AcceptPayment,
			Middlewares: []gin.HandlerFunc{
				middleware.AuthMiddleware(),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/payment/webhook",
			Handler: h.WebHook,
			Middlewares: []gin.HandlerFunc{
				middleware.AuthMiddleware(),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/payments",
			Handler: h.GetPayments,
			Middlewares: []gin.HandlerFunc{
				middleware.AuthMiddleware(),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/payment/:id",
			Handler: h.GetPaymentByCheckoutRequestID,
			Middlewares: []gin.HandlerFunc{
				middleware.AuthMiddleware(),
			},
		},
	}

	route.RegisterRoute(group, paymentRoutes, log)
}
