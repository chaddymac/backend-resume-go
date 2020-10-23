package main

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-lambda-go/events"

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

type Portfolio struct {
	Key string `dynamo:"key" json:"key"` //partitionkey

	VisitorCount int `dynamo:"visitorCount,omitempty" json:"visitorCount"`
}

func Counter() string {
	var dbItem Portfolio
	db := dynamo.New(session.New(), &aws.Config{Region: aws.String("us-east-1")})
	table := db.Table("VisitorCountgo")
	err := table.Update("key", "count").Add("visitorCount", 1).Value(&dbItem)

	if err != nil {
		fmt.Println(err)

	}

	Y := strconv.Itoa(dbItem.VisitorCount)

	return Y
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	x := Counter()
	headers := map[string]string{"Access-Control-Allow-Origin": "*", "Access-Control-Allow-Headers": " Content-Type", "Access-Control-Allow-Methods": "Get"}

	return events.APIGatewayProxyResponse{
		Body:       x,
		StatusCode: 200,
		Headers:    headers,
	}, nil
}

func main() {
	lambda.Start(handler)
}
