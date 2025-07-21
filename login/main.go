package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/golang-jwt/jwt/v5"
)

type Response events.APIGatewayV2HTTPResponse

type Admin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func handler(c context.Context, req events.APIGatewayV2HTTPRequest) (Response, error) {
	var admin Admin
	if err := json.Unmarshal([]byte(req.Body), &admin); err != nil {
		return Response{
			StatusCode: 400,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, err
	}

	if admin.Username != os.Getenv("ADMIN") || admin.Password != os.Getenv("PASS") {
		return Response{
			StatusCode: 401,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, errors.New("unauthorized access")
	}

	claims := jwt.MapClaims{
		"admin": admin.Username,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtKey)

	if err != nil {
		return Response{
			StatusCode: 500,
			Body:       "Token Generation Failed",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, err
	}

	tokenJSON, _ := json.Marshal(map[string]string{"token": tokenStr})

	return Response{
		StatusCode: 200,
		Body:       string(tokenJSON),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
