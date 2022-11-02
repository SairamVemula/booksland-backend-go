package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookList struct {
	BookID  primitive.ObjectID `validate:"required" json:"book_id,omitempty" bson:"book_id,omitempty"`
	StockID primitive.ObjectID `validate:"required" json:"stock_id,omitempty" bson:"stock_id,omitempty"`
}

type Order struct {
	ID              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Type            string             `validate:"required" json:"type,omitempty" bson:"type,omitempty"`
	Items           []BookList         `validate:"required,min=1" json:"items,omitempty" bson:"items,omitempty"`
	Publisher       string             `validate:"required,min=2,max=50" json:"publisher,omitempty" bson:"publisher,omitempty"`
	Year            string             `validate:"required,min=2,max=50" json:"year,omitempty" bson:"year,omitempty"`
	Price           int                `validate:"required,gt=0" json:"price,omitempty" bson:"price,omitempty"`
	DiscountPercent int                `json:"discount_percent,omitempty" bson:"discount_percent,omitempty"`
	DiscountPrice   int64              `json:"discount_price,omitempty" bson:"discount_price,omitempty"`
	CreatedBy       primitive.ObjectID `json:"created_by,omitempty" bson:"created_by,omitempty"`
	PaymentID       string             `json:"payment_id,omitempty" bson:"payment_id,omitempty"`
	PaymentOrderId  string             `json:"payment_order_id,omitempty" bson:"payment_order_id,omitempty"`
	PaymentStatus   string             `json:"status,omitempty" bson:"status,omitempty"`
	Status          string             `json:"payment_status,omitempty" bson:"payment_status,omitempty"`
	PaymentMode     string             `json:"payment_mode,omitempty" bson:"payment_mode,omitempty"`
	CreatedOn       int64              `json:"created_on,omitempty" bson:"created_on,omitempty"`
	UpdatedOn       int64              `json:"updated_on,omitempty" bson:"updated_on,omitempty"`
}

func NewOrder(order *Order) *Order {
	if order.CreatedOn != 0 {
		order.CreatedOn = time.Now().UnixMilli()
	}
	if order.UpdatedOn != 0 {
		order.UpdatedOn = time.Now().UnixMilli()
	}
	return order
}
