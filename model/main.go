package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	FirstName string             `bson:"firstname,omitempty"`
	LastName  string             `bson:"lastname,omitempty"`
	Email     string             `bson:"email,omitempty"`
	Username  string             `bson:"username,omitempty"`
	Password  string             `bson:"password,omitempty"`
}