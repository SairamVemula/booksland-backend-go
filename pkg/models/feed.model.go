package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Feed struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `validate:"required,min=2,max=50" json:"name,omitempty" bson:"name,omitempty"`
	Title     string             `json:"title,omitempty" bson:"title,omitempty"`
	Type      string             `json:"type,omitempty" bson:"type,omitempty"`           // sections,courses,books
	ViewType  string             `json:"view_type,omitempty" bson:"view_type,omitempty"` // banner, 2colgrid, rowgrid,colgrid
	Linked    string             `json:"linked,omitempty" bson:"linked,omitempty"`
	Paralink  string             `validate:"required,min=2,max=50" json:"paralink,omitempty" bson:"paralink,omitempty"`
	Sections  []Section          `json:"sections,omitempty" bson:"sections,omitempty"`
	Order     int                `validate:"required,number" json:"order,omitempty" bson:"order,omitempty"`
	CreatedBy primitive.ObjectID `json:"created_by,omitempty" bson:"created_by,omitempty"`
	CreatedOn int64              `json:"created_on,omitempty" bson:"created_on,omitempty"`
	UpdatedOn int64              `json:"updated_on,omitempty" bson:"updated_on,omitempty"`
}
type Section struct {
	Title    string             `json:"title,omitempty" bson:"title,omitempty"`
	Type     string             `json:"type,omitempty" bson:"type,omitempty"` // paralink,course,book,options
	Paralink string             `json:"paralink,omitempty" bson:"paralink,omitempty"`
	Image    primitive.ObjectID `json:"image,omitempty" bson:"image,omitempty"`
	Course   primitive.ObjectID `json:"course,omitempty" bson:"course,omitempty"`
	Book     primitive.ObjectID `json:"book,omitempty" bson:"book,omitempty"`
	Options  Options            `json:"options,omitempty" bson:"options,omitempty"`
}

type Options struct {
	Title    string    `json:"title,omitempty" bson:"title,omitempty"`
	Sections []Section `json:"sections,omitempty" bson:"sections,omitempty"`
}

func NewFeed(feed *Feed) *Feed {
	if feed.CreatedOn == 0 {
		feed.CreatedOn = time.Now().UnixMilli()
	}
	if feed.UpdatedOn == 0 {
		feed.UpdatedOn = time.Now().UnixMilli()
	}
	return feed
}
