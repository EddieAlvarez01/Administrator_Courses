package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Section struct {
	ID primitive.ObjectID	`bson:"_id,omitempty" json:"id"`
	CourseID primitive.ObjectID	`bson:"course_id" json:"course_id" validate:"required"`
	Professor primitive.ObjectID	`bson:"professor" json:"professor" validate:"required"`
	Name string	`bson:"name" json:"name" validate:"required,alpha"`
	StartHour time.Time	`bson:"start_hour" json:"start_hour" validate:"required"`
	EndHour time.Time	`bson:"end_hour" json:"end_hour" validate:"required"`
}