package main

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Res events.APIGatewayV2HTTPResponse

func handler(ctx context.Context) (Res, error) {
	var buf bytes.Buffer

	body, err := json.Marshal(map[string]interface{}{
		"message": "Hello from Serverless-Go-AWS",
	})
	if err != nil {
		return Res{
			StatusCode: 500,
			Body:       err.Error(),
		}, err
	}
	json.HTMLEscape(&buf, body)

	res := Res{
		StatusCode: 200,
		Body:       buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return res, nil
}

func main() {
	lambda.Start(handler)
}
