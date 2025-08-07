package payment

import (
	"github.com/Kirubel-Enyew27/safari-payment/internal/model/persistencedb"
	"github.com/Kirubel-Enyew27/safari-payment/internal/storage"
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
