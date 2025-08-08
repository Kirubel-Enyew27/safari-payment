package service

import (
	"context"

	"github.com/Kirubel-Enyew27/safari-payment/internal/model/dto"
)

type Payment interface {
	AcceptPayment(ctx context.Context, req dto.AcceptPaymentRequest) (dto.AcceptPaymentResponse, error)
	StorePayment(ctx context.Context, payload map[string]any) (dto.Payment, error)
	GetPayments(ctx context.Context, limit int32, offset int32) ([]dto.Payment, error)
}
