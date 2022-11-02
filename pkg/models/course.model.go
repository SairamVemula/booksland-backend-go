package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Course struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `validate:"required,min=2,max=50" json:"name,omitempty" bson:"name,omitempty"`
	Streams   []string           `json:"streams,omitempty" bson:"streams,omitempty"`
	Semesters []string           `json:"semesters,omitempty" bson:"semesters,omitempty"`
	Tags      []string           `json:"tags,omitempty" bson:"tags,omitempty"`
	Image     primitive.ObjectID `validate:"required" json:"image,omitempty" bson:"image,omitempty"`
	Order     int                `json:"order,omitempty" bson:"order,omitempty"`
	CreatedBy primitive.ObjectID `json:"created_by,omitempty" bson:"created_by,omitempty"`
	CreatedOn int64              `json:"created_on,omitempty" bson:"created_on,omitempty"`
	UpdatedOn int64              `json:"updated_on,omitempty" bson:"updated_on,omitempty"`
}

func NewCourse(course *Course) *Course {
	if course.CreatedOn == 0 {
		course.CreatedOn = time.Now().UnixMilli()
	}
	if course.UpdatedOn == 0 {
		course.UpdatedOn = time.Now().UnixMilli()
	}
	return course
}
