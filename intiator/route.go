package intiator

import (
	"github.com/Kirubel-Enyew27/safari-payment/internal/route/payment"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitRouter(
	group *gin.RouterGroup,
	handler Handler,
	service Service,
	log *zap.Logger,
) {
	payment.InitRoute(group, handler.payment, log)
}
