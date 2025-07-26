package wccrud

import (
	"context"
	"functions/models"
	"functions/shared"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func getOneDailyHandler(c context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	uri := os.Getenv("MONGO_URI")
	var blog models.DailyBlog
	return shared.GetOne(uri, "daily", &blog, "slug", c, req)
}

func main() {
	lambda.Start(getOneDailyHandler)
}
