package payment

import (
	"context"
	"net/http"
	"time"

	"github.com/Kirubel-Enyew27/safari-payment/internal/errors"
	"github.com/Kirubel-Enyew27/safari-payment/internal/handler"
	"github.com/Kirubel-Enyew27/safari-payment/internal/model/dto"
	"github.com/Kirubel-Enyew27/safari-payment/internal/model/response"
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

	response.SendSuccessResponse(c, http.StatusOK, resp, nil)

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

	response.SendSuccessResponse(c, http.StatusOK, resp, nil)
}

func (p *payment) GetPayments(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), p.contextTimeout)
	defer cancel()

	var limit int32 = 10
	var offset int32 = 0

	err := c.ShouldBindQuery(&struct {
		Limit  *int32 `form:"limit"`
		Offset *int32 `form:"offset"`
	}{Limit: &limit, Offset: &offset})

	if err != nil {
		p.logger.Error("invalid query parameters", zap.Error(err))
		_ = c.Error(errors.ErrBadRequest.Wrap(err, "invalid query parameters"))
		return
	}

	payments, err := p.paymentService.GetPayments(ctx, limit, offset)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.SendSuccessResponse(c, http.StatusOK, payments, nil)
}

func (p *payment) GetPaymentByCheckoutRequestID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), p.contextTimeout)
	defer cancel()

	checkout_request_id := c.Param("id")
	if checkout_request_id == "" {
		p.logger.Error("empty url parameter", zap.String("url parameter value", checkout_request_id))
		_ = c.Error(errors.ErrBadRequest.New("empty url parameter"))
		return
	}
	payment, err := p.paymentService.GetPaymentByCheckoutRequestID(ctx, checkout_request_id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.SendSuccessResponse(c, http.StatusOK, payment, nil)

}
