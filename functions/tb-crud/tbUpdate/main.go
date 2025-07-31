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
	var blog models.TechBlog
	return shared.Update(uri, "techblogs", &blog, "slug", c, req)
}

func main() {
	lambda.Start(handler)
}
