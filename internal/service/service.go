package service

import (
	"context"

	"github.com/Kirubel-Enyew27/safari-payment/internal/model/dto"
)

type Payment interface {
	AcceptPayment(ctx context.Context, req dto.AcceptPaymentRequest) (dto.AcceptPaymentResponse, error)
}
