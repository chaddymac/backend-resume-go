package main

import (
	"testing"

	_ "github.com/aws/aws-sdk-go/service/dynamodb"
	_ "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	a := Counter()
	assert.IsType(t, "x", a)

}
