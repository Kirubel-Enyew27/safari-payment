package intiator

import (
	"github.com/Kirubel-Enyew27/safari-payment/internal/model/dto"
	"github.com/Kirubel-Enyew27/safari-payment/internal/service"
	"github.com/Kirubel-Enyew27/safari-payment/internal/service/payment"
	"go.uber.org/zap"
)

type Service struct {
	payment service.Payment
}

func InitService(persistence Persistence, token dto.AccessTokenResponse, log *zap.Logger) Service {

	return Service{
		payment: payment.InitService(persistence.payment, token, log),
	}
}
