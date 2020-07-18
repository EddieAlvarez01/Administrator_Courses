package models

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type assignment struct {
	ID   primitive.ObjectID `bson:"section_id,omitempty" json:"section_id" validate:"required"`
	Date time.Time          `bson:"data,omitempty" json:"date" validate:"required"`
}

//Person .
type Person struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Card        uint64             `bson:"card,omitempty" json:"card" validate:"required,numeric"`
	FirstNames  string             `bson:"first_names,omitempty" json:"first_names" validate:"required"`
	LastNames   string             `bson:"last_names,omitempty" json:"last_names" validate:"required"`
	Address     string             `bson:"address,omitempty" json:"address"`
	Email       string             `bson:"email,omitempty" json:"email" validate:"required,email"`
	Password    string             `bson:"password,omitempty" json:"password,omitempty" validate:"required"`
	Birthdate   time.Time          `bson:"birthdate,omitempty" json:"birthdate" validate:"required"`
	Role        []string           `bson:"role,omitempty" json:"role" validate:"required,ne=0"`
	Phone       uint64             `bson:"phone,omitempty" json:"phone" validate:"numeric"`
	Assignments []assignment       `bson:"assignments" json:"assignments"`
	Token 		string 				`bson:"-" json:"token,omitempty"`
}

//Encrypt (ENCRYPT PASSWORD)
func (p *Person) Encrypt() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(p.Password), 4)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	p.Password = string(bytes)
	return nil
}

//CheckPassword CHECK THAT THE HASH IS EQUAL TO THE PASSWORD
func (p *Person) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.Password), []byte(password))
	return err == nil
}
