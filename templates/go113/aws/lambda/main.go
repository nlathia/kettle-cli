package main
 
import (
        "fmt"
        "github.com/aws/aws-lambda-go/lambda"
)

type RequestEvent struct {
        Name string `json:"What is your name?"`
        Age int     `json:"How old are you?"`
}
 
type ResponseEvent struct {
        Message string `json:"Answer:"`
}

func init() {
	
}
 
func {{.FunctionName}}(event RequestEvent) (ResponseEvent, error) {
        return ResponseEvent{Message: fmt.Sprintf("%s is %d years old!", event.Name, event.Age)}, nil
}
 
func main() {
        lambda.Start({{.FunctionName}})
}
