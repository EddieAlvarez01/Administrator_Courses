package dao

import (
	"context"
	"fmt"
	"github.com/EddieAlvarez01/administrator_courses/dao/mongodb"
	"github.com/EddieAlvarez01/administrator_courses/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type AssignmentImpl struct {}

//GET CLIENT AND COLLECTION
func (dao AssignmentImpl) initDB() (*mongo.Client, *mongo.Collection){
	client := mongodb.GetConnection()
	collection := client.Database("admin_courses").Collection("assignments")
	return client, collection
}

//CREATE A NEW ASSIGNMENT
func (dao AssignmentImpl) Create(assigment *models.Assignment) error {
	client, assignmentsCollection := dao.initDB()
	defer client.Disconnect(context.TODO())
	result, err := assignmentsCollection.InsertOne(context.TODO(), assigment)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = assignmentsCollection.FindOne(context.TODO(), bson.D{{"_id", result.InsertedID}}).Decode(assigment)
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}

//GET ASSIGNMENT BY PERSON ID AND DATE RANGE
func (dao AssignmentImpl) GetByPersonAndDateRange(idPerson string, startDate, endDate time.Time) (*models.Assignment, error) {
	client, assignmentsCollection := dao.initDB()
	defer client.Disconnect(context.TODO())
	objectID, err := primitive.ObjectIDFromHex(idPerson)
	if err != nil {
		return nil, err
	}
	var assigment models.Assignment
	filter := bson.D{
		{"person_id", objectID},
		{"date", bson.M{
			"$gt": startDate,
			"$lt": endDate,
		}},
	}
	err = assignmentsCollection.FindOne(context.TODO(), filter).Decode(&assigment)
	if err != nil {
		return nil, nil
	}
	return &assigment, nil
}
