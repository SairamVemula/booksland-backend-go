package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartItem struct {
	ID              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID          primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	BookID          primitive.ObjectID `validate:"required" json:"book_id,omitempty" bson:"book_id,omitempty"`
	CourseID        primitive.ObjectID `json:"course_id,omitempty" bson:"course_id,omitempty"`
	StockID         primitive.ObjectID `json:"stock_id,omitempty" bson:"stock_id,omitempty"`
	Publisher       string             `validate:"required" json:"publisher,omitempty" bson:"publisher,omitempty"`
	Year            string             `validate:"required" json:"year,omitempty" bson:"year,omitempty"`
	Price           int                `json:"price,omitempty" bson:"price,omitempty"`
	DiscountPercent int                `json:"discount_percent,omitempty" bson:"discount_percent,omitempty"`
	Quantity        int                `json:"quantity,omitempty" bson:"quantity,omitempty"`
	CreatedOn       int64              `json:"created_on,omitempty" bson:"created_on,omitempty"`
	UpdatedOn       int64              `json:"updated_on,omitempty" bson:"updated_on,omitempty"`
}

func NewCartItem(cartitem *CartItem) *CartItem {
	if cartitem.CreatedOn == 0 {
		cartitem.CreatedOn = time.Now().UnixMilli()
	}
	if cartitem.UpdatedOn == 0 {
		cartitem.UpdatedOn = time.Now().UnixMilli()
	}
	if cartitem.Quantity == 0 {
		cartitem.Quantity = 1
	}
	return cartitem
}
