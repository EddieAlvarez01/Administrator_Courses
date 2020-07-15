package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type assignment struct {
	id   primitive.ObjectID `bson:"section_id,omitempty" json:"section_id" validate:"required"`
	date time.Time          `bson:"data,omitempty" json:"date" validate:"required"`
}

//Person .
type Person struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Card        uint64             `bson:"card,omitempty" json:"card" validate:"required,numeric"`
	FirstNames  string             `bson:"first_names,omitempty" json:"first_names" validate:"required"`
	LastNames   string             `bson:"last_names,omitempty" json:"last_names" validate:"required"`
	Address     string             `bson:"address,omitempty" json:"address"`
	Email       string             `bson:"email,omitempty" json:"email" validate:"required,email"`
	Password    string             `bson:"password,omitempty" json:"password" validate:"required"`
	Birthdate   time.Time          `bson:"birthdate,omitempty" json:"birthdate" validate:"required,datetime=2006-01-02"`
	Role        []string           `bson:"role,omitempty" json:"role" validate:"required"`
	Phone       uint64             `bson:"phone,omitempty" json:"phone" validate:"numeric"`
	Assignments []assignment       `bson:"assignments" json:"assignments"`
}
