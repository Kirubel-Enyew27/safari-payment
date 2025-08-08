package storage

import (
	"context"

	"github.com/Kirubel-Enyew27/safari-payment/internal/model/dto"
)

type Payment interface {
	SavePayment(ctx context.Context, payment dto.Payment) (dto.Payment, error)
	GetPayments(ctx context.Context, limit int32, offset int32) ([]dto.Payment, error)
	GetPaymentByCheckoutRequestID(ctx context.Context, id string) (dto.Payment, error)
}
