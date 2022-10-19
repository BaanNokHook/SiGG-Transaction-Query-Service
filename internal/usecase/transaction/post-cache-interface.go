package usecase

import (
	"nextclan/transaction-gateway/transaction-query-service/internal/entity"
)

type QueryRawTxnCache interface {
	//Set(key string, value entity.RawTransaction)
	Get(key string) *entity.RawTransaction
}
