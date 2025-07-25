package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Comment struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string        `json:"name"`
	Comment   string        `json:"comment"`
	CreatedAt time.Time     `json:"createdAt"`
}
