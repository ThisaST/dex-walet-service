package main

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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
			"Type": {
				S: aws.String(wallet.Type),
			},
		},
		TableName: &TableName,
	})

	return err
}

// PutItem inserts the struct Person
func PutItem(person Person) error {
	_, err := dynamo.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"Id": {
				N: aws.String(strconv.Itoa(person.Id)),
			},
			"Name": {
				S: aws.String(person.Name),
			},
			"Website": {
				S: aws.String(person.Website),
			},
		},
		TableName: &TableName,
	})

	return err
}

// UpdateItem updates the Person based on the Person.Id
func UpdateItem(person Person) error {
	_, err := dynamo.UpdateItem(&dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#N": aws.String("Name"),
			"#W": aws.String("Website"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":Name": {
				S: aws.String(person.Name),
			},
			":Website": {
				S: aws.String(person.Website),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				N: aws.String(strconv.Itoa(person.Id)),
			},
		},
		TableName:        &TableName,
		UpdateExpression: aws.String("SET #N = :Name, #W = :Website"),
	})

	return err
}

// DeleteItem deletes a Person based on the Person.Id
func DeleteItem(person Person) error {
	_, err := dynamo.DeleteItem(&dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				N: aws.String(strconv.Itoa(person.Id)),
			},
		},
		TableName: &TableName,
	})

	return err
}

// GetItem gets a Person based on the Id, returns a person
func GetItem(id int) (person Person, err error) {
	result, err := dynamo.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				N: aws.String(strconv.Itoa(id)),
			},
		},
		TableName: &TableName,
	})

	if err != nil {
		return person, err
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &person)

	return person, err

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

func GetItems() (persons *[]Person, err error) {
	result, err := dynamo.Scan(&dynamodb.ScanInput{
		TableName: &TableName,
	})

	if err != nil {
		return persons, err
	}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &persons)

	return persons, err

}
