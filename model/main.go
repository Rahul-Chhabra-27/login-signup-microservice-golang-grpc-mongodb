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

// `bson:"_id,omitempty"` is a tag that tells MongoDB to use the field `_id` as the primary key for the collection.
// `bson:"firstname,omitempty"` is a tag that tells MongoDB to use the field `firstname` as a column in the collection.
// `bson:"lastname,omitempty"` is a tag that tells MongoDB to use the field `lastname` as a column in the collection.
// removing this `bson:"email,omitempty"` tag will cause the email field to be ignored by MongoDB.
