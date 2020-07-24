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

type SectionImpl struct {}

//GET CONNECTION AND COLLECTION OF DB
func (s SectionImpl) initDB() (*mongo.Client, *mongo.Collection){
	client := mongodb.GetConnection()
	db := client.Database("admin_courses")
	return client, db.Collection("sections")
}

//CREATE A NEW SECTION
func (s SectionImpl) Create(section *models.Section) error {
	client, sectionsCollection := s.initDB()
	defer client.Disconnect(context.TODO())
	lastInsert, err := sectionsCollection.InsertOne(context.TODO(), section)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	section.ID = lastInsert.InsertedID.(primitive.ObjectID)
	return nil
}

//GET THE SECTION BY COURSE ID AND SECTION LETTER
func (s SectionImpl) GetByCourseAndSection(idCourse primitive.ObjectID, section string) *models.Section{
	client, sectionsCollection := s.initDB()
	defer client.Disconnect(context.TODO())
	filter := bson.D{
		{"name", section},
		{"course_id", idCourse},
	}
	var findSection models.Section
	err := sectionsCollection.FindOne(context.TODO(), filter).Decode(&findSection)
	if err != nil {
		return nil
	}
	return &findSection
}

//GET THE COURSE PROFESSOR
func (s SectionImpl) GetProfessor(idProfessor primitive.ObjectID) *models.Person {
	client, _ := s.initDB()
	defer client.Disconnect(context.TODO())
	personsCollection := client.Database("admin_courses").Collection("persons")
	filter := bson.D{
		{"_id", idProfessor},
	}
	var person models.Person
	err := personsCollection.FindOne(context.TODO(), filter).Decode(&person)
	if err != nil {
		return nil
	}
	return &person
}

//UPDATE A SECTION
func (s SectionImpl) Update(section *models.Section) error {
	client, sectionsCollection := s.initDB()
	defer client.Disconnect(context.TODO())
	filter := bson.D{
		{"_id", section.ID},
	}
	_, err := sectionsCollection.UpdateOne(context.TODO(), filter, bson.D{
		{"$set", bson.D{
			{"course_id", section.CourseID},
			{"professor", section.Professor},
			{"name", section.Name},
			{"start_hour", section.StartHour},
			{"end_hour", section.EndHour},
		}},
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}

//GET A SECTION BY ID
func (s SectionImpl) GetById(id string) (*models.Section, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	client, sectionsCollection := s.initDB()
	defer client.Disconnect(context.TODO())
	filter := bson.D{
		{"$match", bson.D{
			{"_id", objectID},
		}},
	}
	lookup := bson.D{{"$lookup", bson.D{{"from", "courses"}, {"localField", "course_id"}, {"foreignField", "_id"}, {"as", "course"}}}}
	unwind := bson.D{
		{"$unwind", bson.D{
			{"path", "$course"},
			{"preserveNullAndEmptyArrays", false},
		}},
	}
	var section models.Section
	cursor, err := sectionsCollection.Aggregate(context.TODO(), mongo.Pipeline{filter, lookup, unwind})
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	if cursor == nil {
		return nil, nil
	}
	_ = cursor.Decode(&section)
	return &section, nil
}

//GET ALL SECTIONS BY COURSE ID
func (s SectionImpl) GetAllByCourseID(id primitive.ObjectID) ([]models.Section, error) {
	client, sectionsCollection := s.initDB()
	defer client.Disconnect(context.TODO())
	filter := bson.D{
		{"$match", bson.D{{"course_id", id}}},
	}
	lookup := bson.D{{"$lookup", bson.D{{"from", "courses"}, {"localField", "course_id"}, {"foreignField", "_id"}, {"as", "course"}}}}
	unwind := bson.D{
		{"$unwind", bson.D{
			{"path", "$course"},
			{"preserveNullAndEmptyArrays", false},
		}},
	}
	var sections []models.Section
	cursor, err := sectionsCollection.Aggregate(context.TODO(), mongo.Pipeline{filter, lookup, unwind})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if err = cursor.All(context.TODO(), &sections); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return sections, nil
}
