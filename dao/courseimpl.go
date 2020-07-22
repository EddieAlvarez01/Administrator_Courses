package dao

import (
	"context"
	"fmt"
	"github.com/EddieAlvarez01/administrator_courses/dao/mongodb"
	"github.com/EddieAlvarez01/administrator_courses/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CourseImpl struct {}

//initDB GET A CONNECTION DB AND COURSES COLLECTION
func (c CourseImpl) initDB() (*mongo.Client, *mongo.Collection) {
	client := mongodb.GetConnection()
	db := client.Database("admin_courses")
	collection := db.Collection("courses")
	return client, collection
}

//Create CREATE A NEW COURSE
func (c CourseImpl) Create(course *models.Course) error {
	client, coursesCollection := c.initDB()
	defer client.Disconnect(context.TODO())
	result, err := coursesCollection.InsertOne(context.TODO(), course)
	if err != nil {
		fmt.Print(err.Error())
		return err
	}
	err = coursesCollection.FindOne(context.TODO(), bson.D{{"_id", result.InsertedID}}).Decode(course)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

//GET ONE COURSE
func (c CourseImpl) GetById(id string) (*models.Course, error){
	client, coursesCollection := c.initDB()
	defer client.Disconnect(context.TODO())
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	filter := bson.D{{"_id", objectID}}
	var course models.Course
	err = coursesCollection.FindOne(context.TODO(), filter).Decode(&course)
	if err != nil {
		return nil, nil
	}
	return &course, nil
}

//UPDATE A COURSE
func (c CourseImpl) Update(course *models.Course) error {
	client, coursesCollection := c.initDB()
	defer client.Disconnect(context.TODO())
	filter := bson.D{{"_id", course.ID}}
	_, err := coursesCollection.UpdateOne(context.TODO(), filter, bson.D{
		{"$set", bson.D{
			{"name", course.Name},
			{"credits", course.Credits},
			{"school", course.School},
		}},
	})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

//GET ALL COURSES
func (c CourseImpl) GetAll() ([]*models.Course, error) {
	client, coursesCollection := c.initDB()
	defer client.Disconnect(context.TODO())
	var courses []*models.Course
	cursor, err := coursesCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	if err = cursor.All(context.TODO(), &courses); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return courses, nil
}