package usecase

import (
	"encoding/json"
	"fmt"
	"nextclan/transaction-gateway/transaction-query-service/internal/entity"
	"nextclan/transaction-gateway/transaction-query-service/pkg/logger"

	redisCache "nextclan/transaction-gateway/transaction-query-service/pkg/redis"
)

type QueryRawTransaction struct {
	log         logger.Interface
	redisClient *redisCache.RedisCache
}

func NewQueryRawTxn(l logger.Interface, redisClient *redisCache.RedisCache) *QueryRawTransaction {
	return &QueryRawTransaction{log: l, redisClient: redisClient}
}

func (r *QueryRawTransaction) Get(key string) *entity.RawTransaction {
	val, err := r.redisClient.Get(key)
	rawTransaction := &entity.RawTransaction{}
	err = json.Unmarshal([]byte(val), rawTransaction)
	if err != nil {
		fmt.Println(err)
	}
	return rawTransaction
}
