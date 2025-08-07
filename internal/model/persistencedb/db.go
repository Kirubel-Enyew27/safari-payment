package persistencedb

import (
	"github.com/Kirubel-Enyew27/safari-payment/internal/model/db"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type PersistenceDB struct {
	*db.Queries
	pool *pgxpool.Pool
	log  *zap.Logger
}

func New(pool *pgxpool.Pool, log *zap.Logger) PersistenceDB {
	return PersistenceDB{
		Queries: db.New(pool),
		pool:    pool,
		log:     log,
	}
}
