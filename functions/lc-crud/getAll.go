package lccrud

import (
	"context"
	"functions/models"
	"functions/shared"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func getAllLCHandler(c context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	uri := os.Getenv("MONGO_URI")
	var blogs []models.LCSolution
	return shared.GetAll(uri, "solutions", &blogs, c, req)
}

func main() {
	lambda.Start(getAllLCHandler)
}
