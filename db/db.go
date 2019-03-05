package db

import (
	"os"
	"strings"

	"github.com/tradziej/lambda-ig-followers-chart/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var (
	tableName = os.Getenv("LOGS_TABLE")
	region    = os.Getenv("AWS_REGION")
	localEnv  = os.Getenv("AWS_SAM_LOCAL")
)

type DB struct {
	Instance *dynamodb.DynamoDB
}

func New() *DB {
	awsConfig := &aws.Config{
		Region: aws.String(region),
	}

	// hack for local development
	if len(localEnv) > 0 && strings.ToLower(localEnv) == "true" {
		awsConfig.Endpoint = aws.String("http://host.docker.internal:8000")
	}

	svc := dynamodb.New(
		session.New(),
		awsConfig,
	)
	return &DB{
		Instance: svc,
	}
}

func (db *DB) PutItem(i interface{}) error {
	av, err := dynamodbattribute.MarshalMap(i)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}
	_, err = db.Instance.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetItems() ([]models.DataLog, error) {
	items, err := db.Instance.Query(&dynamodb.QueryInput{
		TableName: aws.String(tableName),
		KeyConditions: map[string]*dynamodb.Condition{
			"username": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(os.Getenv("IG_USERNAME")),
					},
				},
			},
		},
	})

	if err != nil {
		return nil, err
	}
	dataLog := []models.DataLog{}
	err = dynamodbattribute.UnmarshalListOfMaps(items.Items, &dataLog)

	return dataLog, err
}
