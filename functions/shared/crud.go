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
	authHeader := req.Headers["authorization"]
	ok, err := RequireAuth(authHeader)
	if !ok || err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 401,
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

func Update(uri string, collName string, res interface{}, param string, c context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	authHeader := req.Headers["authorization"]
	ok, err := RequireAuth(authHeader)
	if !ok || err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 401,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, err
	}

	slug := req.PathParameters[param]

	if err := json.Unmarshal([]byte(req.Body), res); err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, err
	}

	coll, err := GetCollection(collName, uri)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, err
	}

	val := reflect.ValueOf(res).Elem()
	title := val.FieldByName("Title")
	description := val.FieldByName("Description")
	s := val.FieldByName("Slug")
	content := val.FieldByName("Content")

	update := bson.M{
		"$set": bson.M{
			"title":       title,
			"description": description,
			"slug":        s,
			"content":     content,
		},
	}
	result, e := coll.UpdateOne(c, bson.M{"slug": slug}, update)
	if e != nil || result.MatchedCount == 0 {
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

func Delete(uri string, collName string, res interface{}, param string, c context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	authHeader := req.Headers["authorization"]
	ok, err := RequireAuth(authHeader)
	if !ok || err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 401,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, err
	}

	slug := req.PathParameters[param]

	if err := json.Unmarshal([]byte(req.Body), res); err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, err
	}

	coll, err := GetCollection(collName, uri)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, err
	}

	result, eror := coll.DeleteOne(c, bson.M{"slug": slug})
	if eror != nil || result.DeletedCount == 0 {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, eror
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 204,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}
