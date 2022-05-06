package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"strconv"
)

var dynamo *dynamodb.DynamoDB

// connectDynamo returns a dynamoDB client
func connectDynamo() (db *dynamodb.DynamoDB) {
	return dynamodb.New(session.Must(session.NewSession(&aws.Config{
		Region: &RegionName,
	})))
}

func CreateWalletTable() error {
	_, err := dynamo.CreateTable(&dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Address"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Address"),
				KeyType:       aws.String("HASH"),
			},
		},
		BillingMode: aws.String(dynamodb.BillingModePayPerRequest),
		TableName:   &TableName,
	})

	return err
}

func AddWallet(wallet Wallet) error {
	_, err := dynamo.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"Address": {
				N: aws.String(strconv.Itoa(wallet.Address)),
			},
			"PrivateKey": {
				S: aws.String(wallet.PrivateKey),
			},
			"PublicKey": {
				S: aws.String(wallet.PublicKey),
			},
			"Balance": {
				N: aws.String(strconv.FormatFloat(wallet.Balance, 'E', -1, 32)),
			},
			"CryptoCurrency": {
				S: aws.String(wallet.CryptoCurrency),
			},
		},
		TableName: &TableName,
	})

	return err
}

func GetWallet(address string) (wallet Wallet, err error) {
	result, err := dynamo.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Address": {
				N: aws.String(address),
			},
		},
		TableName: &TableName,
	})

	if err != nil {
		return wallet, err
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &wallet)

	return wallet, err

}
