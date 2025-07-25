package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type TechBlog struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Slug        string        `json:"slug"`
	Content     string        `json:"content"`
	Comments    []Comment     `json:"comments"`
	CreatedAt   time.Time     `json:"createdAt"`
}
