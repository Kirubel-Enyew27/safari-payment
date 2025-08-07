package payment

import (
	"context"
	"net/http"
	"time"

	"github.com/Kirubel-Enyew27/safari-payment/internal/errors"
	"github.com/Kirubel-Enyew27/safari-payment/internal/handler"
	"github.com/Kirubel-Enyew27/safari-payment/internal/model/dto"
	"github.com/Kirubel-Enyew27/safari-payment/internal/model/resonse"
	"github.com/Kirubel-Enyew27/safari-payment/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type payment struct {
	logger         *zap.Logger
	paymentService service.Payment
	contextTimeout time.Duration
}

func InitHandler(paymentService service.Payment, contextTimeout time.Duration,
	log *zap.Logger) handler.Payment {

	return &payment{
		logger:         log,
		paymentService: paymentService,
		contextTimeout: contextTimeout,
	}
}

func (p *payment) AcceptPayment(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, p.contextTimeout)
	defer cancel()

	var req dto.AcceptPaymentRequest
	if err := c.ShouldBind(&req); err != nil {
		err := errors.ErrBadRequest.Wrap(err, "invalid input to accept payment request")
		p.logger.Error("failed to bind accept payment request body", zap.Error(err))
		_ = c.Error(err)
		return
	}

	resp, err := p.paymentService.AcceptPayment(ctx, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resonse.SendSuccessResponse(c, http.StatusOK, resp, nil)

}

func (p *payment) WebHook(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, p.contextTimeout)
	defer cancel()

	var payload map[string]any
	if err := c.BindJSON(&payload); err != nil {
		err := errors.ErrBadRequest.Wrap(err, "invalid ebhook format")
		p.logger.Error("invalid ebhook format", zap.Error(err))
		_ = c.Error(err)
		return
	}
	p.logger.Info("Received Callback--->:", zap.Any("payload", payload))

	resp, err := p.paymentService.StorePayment(ctx, payload)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resonse.SendSuccessResponse(c, http.StatusOK, resp, nil)
}
