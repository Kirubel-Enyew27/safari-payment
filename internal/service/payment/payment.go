package payment

import (
	"github.com/Kirubel-Enyew27/safari-payment/internal/service"
	"github.com/Kirubel-Enyew27/safari-payment/internal/storage"
	"go.uber.org/zap"
)

type payment struct {
	storage storage.Payment
	logger  *zap.Logger
}

func InitService(paymentStorage storage.Payment, log *zap.Logger) service.Payment {
	return &payment{
		storage: paymentStorage,
		logger:  log,
	}

}
