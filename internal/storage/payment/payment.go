package payment

import (
	"context"

	"github.com/Kirubel-Enyew27/safari-payment/internal/errors"
	"github.com/Kirubel-Enyew27/safari-payment/internal/model/db"
	"github.com/Kirubel-Enyew27/safari-payment/internal/model/dto"
	"github.com/Kirubel-Enyew27/safari-payment/internal/model/persistencedb"
	"github.com/Kirubel-Enyew27/safari-payment/internal/storage"
	"github.com/jackc/pgx/v4"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type payment struct {
	db  persistencedb.PersistenceDB
	log *zap.Logger
}

func InitStorage(db persistencedb.PersistenceDB, log *zap.Logger) storage.Payment {
	return &payment{
		db:  db,
		log: log,
	}
}

func (p *payment) SavePayment(ctx context.Context, payment dto.Payment) (dto.Payment, error) {
	arg := db.SavePaymentParams{
		CheckoutRequestID: payment.CheckoutRequestID,
		MerchantRequestID: payment.MerchantRequestID,
		PhoneNumber:       payment.PhoneNumber,
		Amount:            decimal.NewFromFloat(payment.Amount),
		MpesaReceipt:      payment.MpesaReceipt,
		TransactionDate:   payment.TransactionDate,
		ResultCode:        int32(payment.ResultCode),
		ResultDesc:        payment.ResultDesc,
	}

	paymentStorage, err := p.db.Queries.SavePayment(ctx, arg)
	if err != nil {
		p.log.Info("failed to save payment", zap.Error(err))
		return dto.Payment{}, errors.ErrUnableTocreate.Wrap(err, "failed to save payment")
	}

	savedAmount, _ := paymentStorage.Amount.Float64()

	savedPayment := dto.Payment{
		ID:                paymentStorage.ID,
		CheckoutRequestID: paymentStorage.CheckoutRequestID,
		MerchantRequestID: paymentStorage.MerchantRequestID,
		PhoneNumber:       paymentStorage.PhoneNumber,
		Amount:            savedAmount,
		MpesaReceipt:      paymentStorage.MpesaReceipt,
		TransactionDate:   paymentStorage.TransactionDate,
		ResultCode:        int(paymentStorage.ResultCode),
		ResultDesc:        paymentStorage.ResultDesc,
		CreatedAt:         paymentStorage.CreatedAt.Time,
	}

	return savedPayment, nil
}

func (p *payment) GetPayments(ctx context.Context) ([]dto.Payment, error) {
	payments, err := p.db.Queries.ListPayments(ctx)
	if err != nil {
		if err == pgx.ErrNoRows {
			p.log.Error("failed to get payments", zap.Error(err))
			return nil, errors.ErrUnableToGet.Wrap(err, "payments not found")
		}
		p.log.Error("failed to get payments", zap.Error(err))
		return nil, errors.ErrUnableToGet.Wrap(err, "failed to get payments")
	}

	fetchedPayments := make([]dto.Payment, len(payments))

	for i, payment := range payments {
		fetchedAmount, _ := payment.Amount.Float64()
		fetchedPayments[i] = dto.Payment{
			ID:                payment.ID,
			CheckoutRequestID: payment.CheckoutRequestID,
			MerchantRequestID: payment.MerchantRequestID,
			PhoneNumber:       payment.PhoneNumber,
			Amount:            fetchedAmount,
			MpesaReceipt:      payment.MpesaReceipt,
			TransactionDate:   payment.TransactionDate,
			ResultCode:        int(payment.ResultCode),
			ResultDesc:        payment.ResultDesc,
			CreatedAt:         payment.CreatedAt.Time,
		}
	}

	return fetchedPayments, nil
}
