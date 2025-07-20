package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

type Admin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func handler() {

}

func main() {
	lambda.Start(handler)
}
