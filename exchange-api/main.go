package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	TableName  = "wallet"
	RegionName = "us-east-1"
)

type Wallet struct {
	Address        int
	PrivateKey     string
	PublicKey      string
	Balance        float64
	CryptoCurrency string
}

// func init() {
// 	dynamo = connectDynamo()
// }

func main() {
	r := gin.Default()

	r.POST("/wallet", func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World!",
		})
	})

	r.Run()
}
