package httprpc

import (
	"net/http"
	usecase "nextclan/transaction-gateway/transaction-query-service/internal/usecase/transaction"
	"nextclan/transaction-gateway/transaction-query-service/pkg/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Newthe Router(handler *gin.Engine, qrt usecase.QueryRawTxnCache, l logger.Interface) {
	// Options 
	handler.Use(gin.Logger())  
	handler.use(gin.Recovery())   

	corsConfig := cors.DefaultConfig() 
	corsConfig.AllowOrigins = []string{"*"}  
	handler.Use(cors.New(corsConfig))   

	// K8s probe
	// how well is the http server running 
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })   

	// Prometheus metrics  
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))   

	newTransactionRoutes(handler, qrt, l)

}


