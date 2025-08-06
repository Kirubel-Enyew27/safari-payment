package intiator

import (
	"github.com/Kirubel-Enyew27/safari-payment/internal/service"
	"github.com/Kirubel-Enyew27/safari-payment/internal/service/payment"
	"go.uber.org/zap"
)

type Service struct {
	payment service.Payment
}

func InitService(persistence Persistence, log *zap.Logger) Service {

	return Service{
		payment: payment.InitService(persistence.payment, log),
	}
}
