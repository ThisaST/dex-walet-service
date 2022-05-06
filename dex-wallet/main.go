package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type Person struct {
	Id      int
	Name    string
	Website string
}

const (
	DefaultWallet  string = "DEFAULT_WALLET"
	ExchangeWallet        = "EXCHANGE_WALLET"
)

var (
	TableName  = "wallet"
	RegionName = "us-east-1"
)

func init() {
	dynamo = connectDynamo()
}

type Wallet struct {
	Address        int
	PrivateKey     string
	PublicKey      string
	Balance        float64
	CryptoCurrency string
	Type           string
}

func main() {

	r := gin.Default()

	// Create table and add data
	r.GET("/", func(c *gin.Context) {
		err := CreateWalletTable()
		if err != nil {
			log.Println(err)
		}

		wallet := Wallet{
			Address:        123456789,
			PrivateKey:     "8D41627E46D5B8556D0D3E30EC15538E",
			PublicKey:      "19791D9C7D235A1353531B6A9A98098E740F0430",
			Balance:        32000,
			CryptoCurrency: "Bitcoin",
			Type:           DefaultWallet,
		}

		err = AddWallet(wallet)
		if err != nil {
			log.Println(err)
		}
	})

	// Get wallet address
	r.GET("/wallet/:address", func(c *gin.Context) {
		address := c.Param("address")

		wallet, err := GetWallet(address)

		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		add, _ := strconv.Atoi(address)

		if wallet.Address != add {
			c.JSON(http.StatusNotFound, "Not found")
			return
		}

		c.JSON(http.StatusOK, wallet)

	})

	// Create an exchange wallet
	r.POST("/exchange-wallet", func(c *gin.Context) {
		err := CreateWalletTable()
		if err != nil {
			log.Println(err)
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Success",
		})
	})

	// Store funds in wallet (Passing wallet address)
	r.POST("/funds", func(c *gin.Context) {
		err := CreateWalletTable()
		if err != nil {
			log.Println(err)
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Store funds in wallet",
		})
	})

	r.Run()

}
