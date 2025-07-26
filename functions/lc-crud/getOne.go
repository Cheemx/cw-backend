package lccrud

import (
	"context"
	"functions/models"
	"functions/shared"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func getOneLCHandler(c context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	uri := os.Getenv("MONGO_URI")
	var blog models.LCSolution
	return shared.GetOne(uri, "solutions", &blog, "problemNo", c, req)
}

func main() {
	lambda.Start(getOneLCHandler)
}
