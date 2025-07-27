package main

import (
	"context"
	"functions/models"
	"functions/shared"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(c context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	uri := os.Getenv("MONGO_URI")
	var daily models.DailyBlog
	return shared.Create(uri, "daily", &daily, c, req)
}

func main() {
	lambda.Start(handler)
}
