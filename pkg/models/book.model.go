package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string             `validate:"required,min=2,max=50" json:"name,omitempty" bson:"name,omitempty"`
	Publishers []string           `json:"publishers,omitempty" bson:"publishers,omitempty"`
	CourseID   primitive.ObjectID `json:"course_id,omitempty" bson:"course_id,omitempty"`
	Tags       []string           `json:"tags,omitempty" bson:"tags,omitempty"`
	Image      primitive.ObjectID `json:"image,omitempty" bson:"image,omitempty"`
	Order      int                `json:"order,omitempty" bson:"order,omitempty"`
	CreatedBy  primitive.ObjectID `json:"created_by,omitempty" bson:"created_by,omitempty"`
	CreatedOn  int64              `json:"created_on,omitempty" bson:"created_on,omitempty"`
	UpdatedOn  int64              `json:"updated_on,omitempty" bson:"updated_on,omitempty"`
}

func NewBook(book *Book) *Book {
	if book.CreatedOn == 0 {
		book.CreatedOn = time.Now().UnixMilli()
	}
	if book.UpdatedOn == 0 {
		book.UpdatedOn = time.Now().UnixMilli()
	}
	return book
}
