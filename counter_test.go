package main

// snippet-start:[dynamodb.go.update_item.imports]
import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	_ "github.com/aws/aws-sdk-go/aws/awserr"
	_ "github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	_ "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"testing"
)

type mockDynamoDBClient struct {
	dynamodbiface.DynamoDBAPI
}

func (m *mockDynamoDBClient) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	//getting the DynamoDB record based on the key
	var tab, err = svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("VisitorCountgo"),
		Key: map[string]*dynamodb.AttributeValue{
			"key": {
				S: aws.String("count"),
			},
		},
	})
	if err != nil {
		fmt.Println(error(err))

	}
return tab, err
}

func Testcounter(t*testing.T) {

}


func main() {

	// SDK will use to load credentials from the shared credentials file
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	//Create DynamoDB client
	svc := dynamodb.New(sess)
	Testcounter(svc)
//}
