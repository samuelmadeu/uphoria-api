package models

import "gopkg.in/mgo.v2/bson"

type (
	// User represents the structure of our resource
	User struct {
		Id          bson.ObjectId `json:"id" bson:"_id"`
		CompanyName string        `json:"companyname" bson:"companyname"`
		Email       string        `json:"email bson:"email"`
		IsActive    bool          `json:"isActive" bson:"isActive"`
	}

	// All users represented as an array of User
	Users []User
)
