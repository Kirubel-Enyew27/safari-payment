package intiator

import (
	"time"

	"github.com/Kirubel-Enyew27/safari-payment/internal/handler"
	"github.com/Kirubel-Enyew27/safari-payment/internal/handler/payment"
	"go.uber.org/zap"
)

type Handler struct {
	payment handler.Payment
}

func InitHandler(service Service, log *zap.Logger, timeout time.Duration) Handler {

	return Handler{
		payment: payment.InitHandler(service.payment, timeout, log),
	}
}
