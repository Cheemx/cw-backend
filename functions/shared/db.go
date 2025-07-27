package shared

import (
	"sync"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	clientInstance *mongo.Client
	clientOnce     = new(sync.Once)
)

func GetClient(uri string) (*mongo.Client, error) {
	var err error

	clientOnce.Do(func() {
		clientOptions := options.Client().ApplyURI(uri)
		clientInstance, err = mongo.Connect(clientOptions)
		if err != nil {
			println("Mongo connection error:", err.Error())
		}
	})

	return clientInstance, err
}

func GetCollection(name string, uri string) (*mongo.Collection, error) {
	client, err := GetClient(uri)
	if err != nil {
		return nil, err
	}
	return client.Database("blogdb").Collection(name), nil
}
