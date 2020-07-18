package dao

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/EddieAlvarez01/administrator_courses/dao/mongodb"
	"github.com/EddieAlvarez01/administrator_courses/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//PersonImpl .
type PersonImpl struct{}

func (dao PersonImpl) initDb() (*mongo.Client, *mongo.Collection) {
	client := mongodb.GetConnection()
	db := client.Database("admin_courses")
	personsCollection := db.Collection("persons")
	return client, personsCollection
}

//Create (create a new person in the DB)
func (dao PersonImpl) Create(p *models.Person) error {
	client, personsCollection := dao.initDb()
	defer client.Disconnect(context.TODO())
	_, err := personsCollection.InsertOne(context.TODO(), p)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

//Update (UPDATE DATA OF THE USERS)
func (dao PersonImpl) Update(person *models.Person) error {
	client, personsCollection := dao.initDb()
	defer client.Disconnect(context.TODO())
	filter := bson.D{{"_id", person.ID}}
	_, err := personsCollection.UpdateOne(context.TODO(), filter, bson.D{
		{"$set", bson.D{
			{"email", person.Email},
			{"address", person.Address},
			{"phone", person.Phone},
		}},
	})
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	return nil
}

//GetOne (GET A ONE USER OF DB)
func (dao PersonImpl) GetOne(id string) (*models.Person, error) {
	client, personsCollection := dao.initDb()
	defer client.Disconnect(context.TODO())
	var person models.Person
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{"_id", objectID}}
	err = personsCollection.FindOne(context.TODO(), filter).Decode(&person)
	if err != nil {
		return nil, nil
	}
	return &person, nil
}

//GetAllByFilter GET ALL PERSONS BY FILTER AND OPTIONS
func (dao PersonImpl) GetAllByFilter(filter interface{}, opt *options.FindOptions) ([]*models.Person, error) {
	client, personsCollection := dao.initDb()
	defer client.Disconnect(context.TODO())
	cursor, err := personsCollection.Find(context.TODO(), filter, opt)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	var persons []*models.Person
	for cursor.Next(context.TODO()) {
		var person models.Person
		err := cursor.Decode(&person)
		if err != nil {
			log.Fatal(err)
		}
		persons = append(persons, &person)
	}
	return persons, nil
}

//GetByEmail GET PERSON BY EMAIL
func (dao PersonImpl) GetByEmail(email string) *models.Person {
	client, personsCollection := dao.initDb()
	defer client.Disconnect(context.TODO())
	filter := bson.D{{"email", email}}
	var person models.Person
	err := personsCollection.FindOne(context.TODO(), filter).Decode(&person)
	if err != nil {
		return nil
	}
	return &person
}
