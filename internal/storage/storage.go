package storage

import (
	"context"

	"github.com/Kirubel-Enyew27/safari-payment/internal/model/dto"
)

type Payment interface {
	SavePayment(ctx context.Context, payment dto.Payment) (dto.Payment, error)
}
