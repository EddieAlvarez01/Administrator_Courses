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

//UPDATE ASSIGNMENT
func (dao AssignmentImpl) Update(assignment *models.Assignment) error {
	client, assignmentsCollection := dao.initDB()
	defer client.Disconnect(context.TODO())
	filter := bson.D{
		{"_id", assignment.ID},
	}
	update := bson.D{
		{"$set", bson.D{
			{"courses", assignment.Courses},
		}},
	}
	_, err := assignmentsCollection.UpdateOne(context.TODO(), filter, update)
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

//GET ASSIGNMENT BY ID IN A PERIOD
func (dao AssignmentImpl) GetById(assignmentId string) (*models.Assignment, error) {
	client, assignmentsCollection := dao.initDB()
	defer client.Disconnect(context.TODO())
	objectID, err := primitive.ObjectIDFromHex(assignmentId)
	if err != nil {
		return nil, err
	}
	filter := bson.D{
		{"_id", objectID},
	}
	var assignment models.Assignment
	err = assignmentsCollection.FindOne(context.TODO(), filter).Decode(&assignment)
	if err != nil {
		return nil, nil
	}
	return &assignment, nil
}

//GET ALL ASSIGNMENTS IN A PERIOD
func (dao AssignmentImpl) GetAllBySectionIdInAPeriod(sectionID primitive.ObjectID, startDate, endDate time.Time) ([]*models.Assignment, error) {
	client, assignmentsCollection := dao.initDB()
	defer client.Disconnect(context.TODO())
	match := bson.D{
		{"$match", bson.D{
			{"courses", sectionID},
			{"date", bson.D{
				{"$gt", startDate},
				{"$lt", endDate},
			}},
		}},
	}
	lookup := bson.D{
		{"$lookup", bson.D{
			{"from", "persons"},
			{"localField", "person_id"},
			{"foreignField", "_id"},
			{"as", "person"},
		}}}
	unwind := bson.D{
		{"$unwind", bson.D{
			{"path", "$person"},
			{"preserveNullAndEmptyArrays", false},
		}},
	}
	var assignments []*models.Assignment
	cursor, err := assignmentsCollection.Aggregate(context.TODO(), mongo.Pipeline{match, lookup, unwind})
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	if err = cursor.All(context.TODO(), &assignments); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return assignments, nil
}
