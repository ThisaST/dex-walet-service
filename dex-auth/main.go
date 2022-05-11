package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	TableName  = "user"
	RegionName = "us-east-1"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	awsKeyId := os.Getenv("AWS_ACCESS_KEY_ID")
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")

	dynamo = connectDynamo(awsKeyId, awsAccessKey, region)
}

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
}

type Response struct {
	token string
	user  User
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	appAddr := ":" + os.Getenv("PORT")
	var router = gin.Default()

	err := CreateUserTable()
	if err != nil {
		log.Println(err)
	}

	user := User{
		Id:       123456789,
		Name:     "Thisara",
		Email:    "thisa1@gmail.com",
		Password: "user@123",
	}

	err = Register(user)
	if err != nil {
		log.Println(err)
	}

	router.POST("/login", func(c *gin.Context) {
		body := User{}

		// using BindJson method to serialize body with struct
		if err := c.BindJSON(&body); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		user, err := Login(body)

		fmt.Println(err)
		if err != nil {
			c.JSON(http.StatusUnauthorized, err)
			return
		}

		res := Response{
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			user:  user,
		}

		c.JSON(http.StatusOK, gin.H{
			"token": res.token,
			"user":  res.user,
		})
	})

	srv := &http.Server{
		Addr:    appAddr,
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	//Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
