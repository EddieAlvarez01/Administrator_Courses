package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Course struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name string `bson:"name" json:"name" validate:"required"`
	Credits int `bson:"credits" json:"credits" validate:"required,numeric"`
	School string `bson:"school" json:"school" validate:"required"`
}
