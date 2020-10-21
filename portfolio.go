package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	_ "github.com/aws/aws-sdk-go/aws/awserr"
	_ "github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	_ "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/guregu/dynamo"
	_ "github.com/guregu/dynamo"
)

type portfolio struct {
	key string `dynamo:"key"` //partition key

	visitorCount int `dynamo:"visitorCount,omitempty"`
}

func handler() error {
	var dbItem portfolio
	db := dynamo.New(session.New(), &aws.Config{Region: aws.String("us-east-1")})
	table := db.Table("VisitorCountgo")
	err := table.Update("key", "count").Add("visitorCount", 1).Value(&dbItem)
	if err != nil {
		return err
	}
	return nil
}
func main() {
	lambda.Start(handler)
}
