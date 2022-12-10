package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Verified struct {
	Email bool `json:"email" bson:"email"`
	Phone bool `json:"phone" bson:"phone"`
}

type User struct {
	ID                 primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Type               string             `json:"type,omitempty" bson:"type,omitempty"`
	Name               string             `validate:"required,min=3,max=40" json:"name,omitempty" bson:"name,omitempty"`
	Phone              string             `validate:"required,numeric,min=10,max=10" json:"phone,omitempty" bson:"phone,omitempty"`
	Email              string             `validate:"required,email" json:"email,omitempty" bson:"email,omitempty"`
	Password           string             `validate:"required,min=8,max=20,passwd" json:"password,omitempty" bson:"password,omitempty"`
	Location           []string           `json:"location,omitempty" bson:"location,omitempty"`
	Address            string             `json:"address,omitempty" bson:"address,omitempty"`
	City               string             `json:"city,omitempty" bson:"city,omitempty"`
	State              string             `json:"state,omitempty" bson:"state,omitempty"`
	Pincode            string             `json:"pincode,omitempty" bson:"pincode,omitempty"`
	Token              string             `json:"token,omitempty" bson:"token,omitempty"`
	TokenExpiry        int64              `json:"token_expiry,omitempty" bson:"token_expiry,omitempty"`
	RefreshToken       string             `json:"refresh_token,omitempty" bson:"refresh_token,omitempty"`
	RefreshTokenExpiry int64              `json:"refresh_token_expiry,omitempty" bson:"refresh_token_expiry,omitempty"`
	Verified           Verified           `json:"verified,omitempty" bson:"verified,omitempty"`
	CreatedOn          int64              `json:"created_on,omitempty" bson:"created_on,omitempty"`
	UpdatedOn          int64              `json:"updated_on,omitempty" bson:"updated_on,omitempty"`
}

func NewUser(user *User) {
	user.Verified = Verified{false, false}
	if user.CreatedOn != 0 {
		user.CreatedOn = time.Now().UnixMilli()
	}
	if user.UpdatedOn != 0 {
		user.UpdatedOn = time.Now().UnixMilli()
	}
	if user.Type != "" {
		user.Type = "user"
	}
	return
}

type UpdateUser struct {
	Name         string   `validate:"min_len=3,max_len=40" json:"name,omitempty" bson:"name,omitempty"`
	Phone        string   `validate:"min_len=10,max_len=10,regexp=^[0-9]*$" json:"phone,omitempty" bson:"phone,omitempty"`
	Email        string   `json:"email,omitempty" bson:"email,omitempty"`
	Location     []string `validate:"location" json:"location,omitempty" bson:"location,omitempty"`
	Address      string   `json:"address,omitempty" bson:"address,omitempty"`
	Verified     Verified `json:"verified,omitempty" bson:"verified,omitempty"`
	RefreshToken string   `json:"refresh_token,omitempty" bson:"refresh_token,omitempty"`
	UpdatedOn    int      `json:"updated_on,omitempty" bson:"updated_on,omitempty"`
}

type ChangePassword struct {
	Password    string `validate:"required" json:"password,omitempty" bson:"password,omitempty"`
	NewPassword string `validate:"required,min=8,max=20,passwd" json:"new_password,omitempty" bson:"new_password,omitempty"`
}

type LoginUser struct {
	Username string `validate:"required" json:"username,omitempty" bson:"username,omitempty"`
	Password string `validate:"required,min=8,max=20,passwd" json:"password,omitempty" bson:"password,omitempty"`
}
