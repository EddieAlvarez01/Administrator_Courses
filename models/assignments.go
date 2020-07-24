package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Assignment struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PersonID primitive.ObjectID	`bson:"person_id" json:"person_id,omitempty"`
	Person Person	`bson:"person,omitempty" json:"person,omitempty"`
	Courses []primitive.ObjectID	`bson:"courses" json:"courses" validate:"required,ne=0"`
	Date time.Time          `bson:"date,omitempty" json:"date"`
}
