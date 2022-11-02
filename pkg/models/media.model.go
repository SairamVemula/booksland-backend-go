package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Media struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Path      string             `json:"path,omitempty" bson:"path,omitempty"`
	CreatedBy primitive.ObjectID `json:"created_by,omitempty" bson:"created_by,omitempty"`
	CreatedOn int64              `json:"created_on,omitempty" bson:"created_on,omitempty"`
	UpdatedOn int64              `json:"updated_on,omitempty" bson:"updated_on,omitempty"`
}

func NewMedia(media *Media) *Media {
	if media.CreatedOn == 0 {
		media.CreatedOn = time.Now().UnixMilli()
	}
	if media.UpdatedOn == 0 {
		media.UpdatedOn = time.Now().UnixMilli()
	}
	return media
}
