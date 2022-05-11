package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

var (
	TableName  = "wallet"
	RegionName = "us-east-1"
)

type Request struct {
	Address       int
	StartCurrency string
	ToCurrency    string
	Value         float64
}

type Wallet struct {
	Address        int
	PrivateKey     string
	PublicKey      string
	Balance        float64
	CryptoCurrency string
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	awsKeyId := os.Getenv("AWS_ACCESS_KEY_ID")
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")

	fmt.Println(awsAccessKey, awsKeyId, region)

	dynamo = connectDynamo(awsKeyId, awsAccessKey, region)
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.POST("/exchange", func(c *gin.Context) {

		body := Request{}

		// using BindJson method to serialize body with struct
		if err := c.BindJSON(&body); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		wallet, err := GetWallet(body.Address)

		if err != nil {
			c.JSON(http.StatusNotFound, err)
			return
		}

		// Exchange conversion here
		wallet.Balance = wallet.Balance + body.Value

		err = UpdateWallet(wallet)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"wallet": wallet,
		})
	})

	r.Run(":3000")
}
