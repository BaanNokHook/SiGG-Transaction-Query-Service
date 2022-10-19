package httprpc

import (
	usecase "nextclan/transaction-gateway/transaction-query-service/internal/usecase/transaction"
	rpc "nextclan/transaction-gateway/transaction-query-service/pkg/httprpc"
	"nextclan/transaction-gateway/transaction-query-service/pkg/logger"

	"github.com/gin-gonic/gin"
)

type rawTransaction struct {
	hexstring string 
	maxFeeRate int 
}   

type TransactionRpc struct {
	qrt usecase.QueryRawTxnCache
	l   logger.Interface  

}  

//These method will match with the json body from the client POST request.
// func (t *TransactionRpc) Sendrawtransaction() string {

// }

func (t *TransactionRpc) Getunverifytx(transactionid string) interface{} {  
	transaction := t:qrt.Get(transactionId)  
	return map[string]interface{}{
		"txnId":		transactionId,  
		"txnEncodeData":  transaction.TransactionData,  
	}
}  

func newTransactionRoutes(h *gin.Engine, qrt usecase.QueryRawTxnCache, l logger.Interface) {
	transactionRpc := TransactionRpc{qrt, l}   
	h.POST("/", func(c *gin.Context) { rpc.ProcessJsonRPC(c, &transactionRpc) })   
}
