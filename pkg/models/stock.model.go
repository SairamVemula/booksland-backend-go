package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Stock struct {
	ID              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	BookID          primitive.ObjectID `validate:"required" json:"book_id,omitempty" bson:"book_id,omitempty"`
	CourseID        primitive.ObjectID `json:"course_id,omitempty" bson:"course_id,omitempty"`
	Publisher       string             `validate:"required,min=2,max=50" json:"publisher,omitempty" bson:"publisher,omitempty"`
	Year            string             `validate:"required,min=4,max=4" json:"year,omitempty" bson:"year,omitempty"`
	Price           int                `validate:"required,gt=0" json:"price,omitempty" bson:"price,omitempty"`
	DiscountPercent int                `json:"discount_percent,omitempty" bson:"discount_percent,omitempty"`
	Status          string             `json:"status,omitempty" bson:"status,omitempty"` //available,sold
	CreatedBy       primitive.ObjectID `json:"created_by,omitempty" bson:"created_by,omitempty"`
	CreatedOn       int64              `json:"created_on,omitempty" bson:"created_on,omitempty"`
	UpdatedOn       int64              `json:"updated_on,omitempty" bson:"updated_on,omitempty"`
}

func NewStock(stock *Stock) *Stock {
	if stock.CreatedOn == 0 {
		stock.CreatedOn = time.Now().UnixMilli()
	}
	if stock.UpdatedOn == 0 {
		stock.UpdatedOn = time.Now().UnixMilli()
	}
	if stock.Status == "" {
		stock.Status = "available"
	}
	return stock
}
