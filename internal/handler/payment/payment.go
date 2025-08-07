package payment

import (
	"time"

	"github.com/Kirubel-Enyew27/safari-payment/internal/handler"
	"github.com/Kirubel-Enyew27/safari-payment/internal/service"
	"go.uber.org/zap"
)

type customer struct {
	logger         *zap.Logger
	paymentService service.Payment
	contextTimeout time.Duration
}

func InitHandler(paymentService service.Payment, contextTimeout time.Duration,
	log *zap.Logger) handler.Payment {

	return &customer{
		logger:         log,
		paymentService: paymentService,
		contextTimeout: contextTimeout,
	}
}
