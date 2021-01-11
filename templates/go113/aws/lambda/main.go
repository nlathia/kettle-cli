package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

type RequestEvent struct{}

type ResponseEvent struct {
	Message string `json:"prediction:"`
}

func init() {}

func LambdaHandler(event RequestEvent) (ResponseEvent, error) {
	return ResponseEvent{
		Message: "hello world!",
	}, nil
}

func main() {
	lambda.Start(LambdaHandler)
}
