package payment

import (
	"github.com/Kirubel-Enyew27/safari-payment/internal/model/dto"
	"github.com/Kirubel-Enyew27/safari-payment/internal/service"
	"github.com/Kirubel-Enyew27/safari-payment/internal/storage"
	"go.uber.org/zap"
)

type payment struct {
	storage storage.Payment
	token   dto.AccessTokenResponse
	logger  *zap.Logger
}

func InitService(paymentStorage storage.Payment, token dto.AccessTokenResponse, log *zap.Logger) service.Payment {
	return &payment{
		storage: paymentStorage,
		token:   token,
		logger:  log,
	}

}
