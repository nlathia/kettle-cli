package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type RequestEvent struct {
	Name string `json:"name"`
}

type ResponseEvent struct {
	Message string `json:"answer:"`
}

func init() {
}

func LambdaHandler(event RequestEvent) (ResponseEvent, error) {
	return ResponseEvent{
		Message: fmt.Sprintf("Hello, %s", event.Name),
	}, nil
}

func main() {
	lambda.Start(LambdaHandler)
}
