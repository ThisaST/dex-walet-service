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

func connectDynamo(awsKeyId string, awsAccessKey string, region string) (db *dynamodb.DynamoDB) {

	fmt.Println("codes", awsAccessKey, awsKeyId, region)

	return dynamodb.New(session.Must(session.NewSession(&aws.Config{
		Region:      &RegionName,
		Credentials: credentials.NewStaticCredentials(awsKeyId, awsAccessKey, ""),
	})))
}

func CreateUserTable() error {
	_, err := dynamo.CreateTable(&dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Email"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Email"),
				KeyType:       aws.String("HASH"),
			},
		},
		BillingMode: aws.String(dynamodb.BillingModePayPerRequest),
		TableName:   &TableName,
	})

	return err
}

func Register(user User) error {
	_, err := dynamo.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"Email": {
				S: aws.String(user.Email),
			},
			"Id": {
				N: aws.String(strconv.Itoa(user.Id)),
			},
			"Name": {
				S: aws.String(user.Name),
			},
			"Password": {
				S: aws.String(user.Password),
			},
		},
		TableName: &TableName,
	})

	return err
}

func Login(user User) (authUser User, err error) {
	result, err := dynamo.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Email": {
				S: aws.String(user.Email),
			},
		},
		TableName: &TableName,
	})

	if err != nil {
		return authUser, err
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &authUser)

	return authUser, err

}
