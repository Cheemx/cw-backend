package shared

import (
	"context"
	"encoding/json"
	"reflect"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func GetAll(uri string, collName string, res interface{}, c context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	coll, e := GetCollection(collName, uri)
	if e != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, e
	}

	cursor, err := coll.Find(c, bson.M{})
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, err
	}

	if erro := cursor.All(c, res); erro != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, erro
	}

	body, prob := json.Marshal(res)
	if prob != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, prob
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil

}

func GetOne(uri string, collName string, res interface{}, param string, c context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	slug := req.PathParameters[param]
	coll, err := GetCollection(collName, uri)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, err
	}

	if e := coll.FindOne(c, bson.M{param: slug}).Decode(res); e != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 404,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, e
	}

	body, prob := json.Marshal(res)
	if prob != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, prob
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func Create(uri string, collName string, res interface{}, c context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	authHeader := req.Headers["Authorization"]
	ok, err := RequireAuth(authHeader)
	if !ok || err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, err
	}

	if err := json.Unmarshal([]byte(req.Body), res); err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, err
	}

	val := reflect.ValueOf(res).Elem()
	field := val.FieldByName("CreatedAt")
	if field.IsValid() && field.CanSet() && field.Type() == reflect.TypeOf(time.Time{}) {
		field.Set(reflect.ValueOf(time.Now()))
	}

	coll, erro := GetCollection(collName, uri)
	if erro != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, err
	}

	_, e := coll.InsertOne(c, res)
	if e != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, e
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 201,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}
