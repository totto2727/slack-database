package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sfn"
)

type MyEvent struct {
	Channel     string      `json:"channel"`
	API         string      `json:"api"`
	TwitterList string      `json:"twitter_list"`
	UpdatedAt   string      `json:"updated_at"`
	MaxResults  json.Number `json:"max_results"`
	URL         string      `json:"url"`
}

func handler(ctx context.Context, event events.SNSEvent) {
	mySession := session.Must(session.NewSession())
	stepFunction := sfn.New(mySession)
	stepFunctionArn := os.Getenv("STEP_FUNCTION_ARN")

	for _, record := range event.Records {
		log.Printf("%#v", record.SNS.Message)

		myEvent := MyEvent{}
		err := json.Unmarshal([]byte(record.SNS.Message), &myEvent)
		if err != nil {
			log.Panicln(err)
		}
		log.Printf("%#v", myEvent)

		myEventJson, err := json.Marshal(myEvent)
		if err != nil {
			log.Panicln(err)
		}
		myEventJsonString := string(myEventJson)
		log.Printf("%#v", myEventJsonString)

		log.Println("Execute Step Function")
		output, err := stepFunction.StartExecutionWithContext(ctx, &sfn.StartExecutionInput{Input: &myEventJsonString, StateMachineArn: &stepFunctionArn})
		if err != nil {
			log.Panicln(err)
		}
		log.Println(output)
	}
}

func main() {
	lambda.Start(handler)
}
