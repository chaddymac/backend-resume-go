package main

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"testing"
)

// Define stub
type mockDynamoDBClient struct {
	dynamodbiface.DynamoDBAPI
}

func (m *mockDynamoDBClient) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	// Make response
	key := dynamodb.AttributeValue{}
	key.SetS("count")
	val := dynamodb.AttributeValue{}
	val.SetN("0")
	resp := make(map[string]*dynamodb.AttributeValue)
	resp["key"] = &key
	resp["visitorCount"] = &val

	// Returned canned response
	output := &dynamodb.GetItemOutput{
		Item: resp,
	}
	return output, nil
}
func (m *mockDynamoDBClient) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	// mock response/functionality
	finOut := &dynamodb.PutItemOutput{}
	return finOut,nil
}


// Sample Test Case
func TestDynamodb(t *testing.T) {
	mockSvc := &mockDynamoDBClient{}
	Counter(mockSvc)
}
