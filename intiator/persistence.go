package intiator

import (
	"github.com/Kirubel-Enyew27/safari-payment/internal/model/persistencedb"
	"github.com/Kirubel-Enyew27/safari-payment/internal/storage"
	"github.com/Kirubel-Enyew27/safari-payment/internal/storage/payment"
	"go.uber.org/zap"
)

type Persistence struct {
	payment storage.Payment
}

func InitPersistence(db persistencedb.PersistenceDB, log *zap.Logger) Persistence {
	return Persistence{
		payment: payment.InitStorage(db, log),
	}
}
