package main

// snippet-start:[dynamodb.go.update_item.imports]
import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	_ "github.com/aws/aws-sdk-go/aws/awserr"
	_ "github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	_ "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"strconv"
)


func Counter(svc dynamodbiface.DynamoDBAPI) string {
	//getting the DynamoDB record based on the key
	var result, err = svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("VisitorCountgo"),
		Key: map[string]*dynamodb.AttributeValue{
			"key": {
				S: aws.String("count"),
			},
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	fmt.Println(result)
	fmt.Printf("%T\n", result)

	//unmarshalls the visitorcount of the record above so we can use that value
	count := 0
	err = dynamodbattribute.Unmarshal(result.Item["visitorCount"], &count)

	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	fmt.Println(count)
	fmt.Printf("%T\n", count)

	//adding 1 to the dynamo value
	count = count + 1

	//converting INT dynamo value to a string
	stCount := strconv.Itoa(count)

	fmt.Println(count)
	//
	////Put new count in the table VisitorCount
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"key": {
				S: aws.String("count"),
			}, "visitorCount":{
				N: aws.String(stCount),
			},
		},
		TableName: aws.String("VisitorCountgo"),
	}
	resultPut, err := svc.PutItem(input)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("final",resultPut)
	fmt.Println("Hi! You did it Cha'Diamond")
	return stCount

}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse,error) {
	// SDK will use to load credentials from the shared credentials file
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	//Create DynamoDB client
	svc := dynamodb.New(sess)

	y := Counter(svc)
	headers := map[string]string{"Access-Control-Allow-Origin": "*", "Access-Control-Allow-Headers": " Content-Type","Access-Control-Allow-Methods":"Get"}

	return events.APIGatewayProxyResponse{
		Body:       y, 
		StatusCode: 200,
		Headers: headers,
	},nil
}

func main(){
	lambda.Start(handler)
}



