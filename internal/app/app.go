package app

import (
	"fmt"
	"nextclan/transaction-gateway/transaction-query-service/config"
	"nextclan/transaction-gateway/transaction-query-service/internal/controller/httprpc"
	usecase "nextclan/transaction-gateway/transaction-query-service/internal/usecase/transaction"
	rpc "nextclan/transaction-gateway/transaction-query-service/pkg/httprpc"
	"nextclan/transaction-gateway/transaction-query-service/pkg/logger"
	"nextclan/transaction-gateway/transaction-query-service/pkg/redis"
	"os"
	"os/signal"
	"syscall"

	redisCache "nextclan/transaction-gateway/transaction-query-service/pkg/redis"

	"github.com/gin-gonic/gin"
)


func Run(cfg *config.Config) {

	l := logger.New(cfg.Log.Level)  
	fmt.Println("Starting App...")   
	
	//usecase 
	redisClient := initialRedisClient(cfg)   
	queryRawTransactionUseCase := usecase.NewQueryRawTxn(l, redisClient)   
	httpServer := initializeRPC(l, httpServer, redisClient)  
}

func initializeRPC(l *logger.Logger, qrt *usecase.QueryRawTransaction, cfg *config.Config) *rpc.Server {
	handler := gin.New()  
	httprpc.NewRouter(handler, qrt, l)   
	httpServer := rpc.New(handler, rpc.port(cfg.HTTP.Port))     
	return httpServer  
}  


func ShutdownApplicationHandler(l *logger.Logger, httpServer *rpc.Server, redisClient *redis.RedisCache) {
	// Waiting signal  
	interrupt := make(chan os.Signal, 1)   
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)   

	select {
	case s := <-interrupt:  
		l.Info("app - Run - signal: " + s.String())   
	}  

	err := httpServer.Shutdown()   
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))   
	}   

	err = redisClient.Close()   
	if err != nil {
		l.Error(fmt.Errorf("app - Run - RedisClient.Shutdown: %w", err))   
	}    
}   
