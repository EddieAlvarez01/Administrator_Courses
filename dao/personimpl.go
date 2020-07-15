package dao

import (
	"context"
	"fmt"

	"github.com/EddieAlvarez01/administrator_courses/dao/mongodb"
	"github.com/EddieAlvarez01/administrator_courses/models"
)

//PersonImpl .
type PersonImpl struct{}

//Create (create a new person in the DB)
func (dao PersonImpl) Create(p *models.Person) error {
	client := mongodb.GetConnection()
	defer client.Disconnect(context.TODO())
	db := client.Database("admin_courses")
	personsCollection := db.Collection("persons")
	_, err := personsCollection.InsertOne(context.TODO(), p)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
