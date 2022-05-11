package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"strconv"
)

var dynamo *dynamodb.DynamoDB

// connectDynamo returns a dynamoDB client
func connectDynamo(awsKeyId string, awsAccessKey string, region string) (db *dynamodb.DynamoDB) {

	fmt.Println("codes", awsAccessKey, awsKeyId, region)

	return dynamodb.New(session.Must(session.NewSession(&aws.Config{
		Region:      &RegionName,
		Credentials: credentials.NewStaticCredentials(awsKeyId, awsAccessKey, ""),
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

func UpdateWallet(wallet Wallet) error {
	_, err := dynamo.UpdateItem(&dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#N": aws.String("Balance"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":Balance": {
				N: aws.String(strconv.FormatFloat(wallet.Balance, 'E', -1, 32)),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"Address": {
				N: aws.String(strconv.Itoa(wallet.Address)),
			},
		},
		TableName:        &TableName,
		UpdateExpression: aws.String("SET #N = :Balance"),
	})

	return err
}

func GetWallet(address int) (wallet Wallet, err error) {
	result, err := dynamo.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Address": {
				N: aws.String(strconv.Itoa(address)),
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
