package handlers

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/EddieAlvarez01/administrator_courses/models"

	"github.com/EddieAlvarez01/administrator_courses/authorization"
	"github.com/EddieAlvarez01/administrator_courses/dao/interfaces"

	"github.com/go-playground/validator/v10"

	"github.com/gorilla/mux"
)

type personHandler struct {
	Persondao interfaces.PersonDao
}

//NewPersonHandler return a new instance of NewPersonHandler
func NewPersonHandler(person interfaces.PersonDao) personHandler {
	return personHandler{person}
}

//CreatePerson .
func (p personHandler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		models.NewResponseJSON(w, http.StatusBadRequest, "Send a valid information", nil)
		return
	}
	var person models.Person
	err = json.Unmarshal(req, &person)
	if err != nil {
		models.NewResponseJSON(w, http.StatusInternalServerError, fmt.Sprintf("Server error: %s", err.Error()), nil)
		return
	}

	//VALIDATE DATA
	validate := validator.New()
	err = validate.Struct(person)
	if err != nil {
		models.NewResponseJSON(w, http.StatusBadRequest, fmt.Sprintf("Invalid data: %+v", err.(validator.ValidationErrors)), nil)
		return
	}

	//verify that the mail and the card are unique
	filter := bson.M{"$or": []bson.M{bson.M{"email": person.Email}, bson.M{"card": person.Card}}}
	persons, err := p.Persondao.GetAllByFilter(filter, nil)
	if err != nil {
		models.NewResponseJSON(w, http.StatusInternalServerError, "Error on server", nil)
		return
	}
	if len(persons) > 0 {
		models.NewResponseJSON(w, http.StatusBadRequest, "The card or email already exist", nil)
		return
	}

	//ENCRYPT PASSWORD
	err = person.Encrypt()
	if err != nil {
		models.NewResponseJSON(w, http.StatusInternalServerError, "Error on server", nil)
		return
	}
	err = p.Persondao.Create(&person)
	if err != nil {
		models.NewResponseJSON(w, http.StatusInternalServerError, fmt.Sprintf("Server error: %s", err.Error()), nil)
		return
	}
	models.NewResponseJSON(w, http.StatusCreated, "Person created successfully", nil)
}

//GetOne .
func (p personHandler) GetOne(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	person, err := p.Persondao.GetOne(id)
	person.Password = ""
	if err != nil {
		models.NewResponseJSON(w, http.StatusBadRequest, "ID invalid", nil)
		return
	}
	if person == nil {
		models.NewResponseJSON(w, http.StatusNotFound, "Person not found", nil)
		return
	}
	models.NewResponseJSON(w, http.StatusOK, "OK", person)
}

//SignIn LOGIN IN THE SYSTEM
func (p personHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		models.NewResponseJSON(w, http.StatusBadRequest, "Insert a valid information", nil)
		return
	}
	var person models.Person
	json.Unmarshal(req, &person)
	validate := validator.New()
	err = validate.StructPartial(person, "Email", "Password")
	if err != nil {
		models.NewResponseJSON(w, http.StatusBadRequest, fmt.Sprintf("Invalid data: %+v", err.(validator.ValidationErrors)), nil)
		return
	}
	findPerson := p.Persondao.GetByEmail(person.Email)
	if findPerson == nil {
		models.NewResponseJSON(w, http.StatusBadRequest, "Email incorrect", nil)
		return
	}
	if !findPerson.CheckPassword(person.Password) {
		models.NewResponseJSON(w, http.StatusBadRequest, "Password incorrect", nil)
		return
	}
	token, err := authorization.GenerateToken(*findPerson)
	if err != nil {
		models.NewResponseJSON(w, http.StatusInternalServerError, "Error when generating token", nil)
		return
	}
	findPerson.Token = token
	findPerson.Password = ""
	models.NewResponseJSON(w, http.StatusOK, "OK", findPerson)
}

//Update (UPDATE A ACCOUNT PERSON)
func (p personHandler) Update(w http.ResponseWriter, r *http.Request) {
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		models.NewResponseJSON(w, http.StatusBadRequest, "Invalid json", nil)
		return
	}
	var personToUpdate models.Person
	json.Unmarshal(req, &personToUpdate)
	validate := validator.New()
	err = validate.StructPartial(personToUpdate, "Email", "Phone", "Address")
	payload := r.Context().Value("payload").(models.Payload)
	person, err := p.Persondao.GetOne(payload.ID)
	if err != nil {
		models.NewResponseJSON(w, http.StatusBadRequest, "Invalid id", nil)
		return
	}
	if person == nil {
		models.NewResponseJSON(w, http.StatusNotFound, "Person not found", nil)
		return
	}
	if person.Email != personToUpdate.Email {
		if !p.verifyEmail(personToUpdate.Email) {
			models.NewResponseJSON(w, http.StatusBadRequest, "This email already taken", nil)
			return
		}
		person.Email = personToUpdate.Email
	}
	person.Address = personToUpdate.Address
	person.Phone = personToUpdate.Phone
	err = p.Persondao.Update(person)
	if err != nil {
		models.NewResponseJSON(w, http.StatusInternalServerError, "Error when updating a person", nil)
		return
	}
	person.Password = ""
	models.NewResponseJSON(w, http.StatusOK, "OK", person)
}

//CREATE A NEW PROFESSOR
func (p personHandler) CreateProfessor(w http.ResponseWriter, r *http.Request) {
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		models.NewResponseJSON(w, http.StatusBadRequest, "Invalid json", nil)
		return
	}
	var person models.Person
	err = json.Unmarshal(req, &person)
	if err != nil {
		models.NewResponseJSON(w, http.StatusBadRequest, "Invalid json", nil)
		return
	}
	validate := validator.New()
	err = validate.StructExcept(person, "Card", "Role")
	if err != nil {
		models.NewResponseJSON(w, http.StatusBadRequest, fmt.Sprintf("Errors: %v", err.(validator.ValidationErrors)), nil)
		return
	}
	if !p.verifyEmail(person.Email) {
		models.NewResponseJSON(w, http.StatusBadRequest, "The email already taken", nil)
		return
	}
	person.Card = 0
	person.Role = []string{"PROFESSOR"}
	err = person.Encrypt()
	if err != nil {
		models.NewResponseJSON(w, http.StatusInternalServerError, "Error in server", nil)
		return
	}
	err = p.Persondao.CreateProfessor(&person)
	if err != nil {
		models.NewResponseJSON(w, http.StatusInternalServerError, "Error in server", nil)
		return
	}
	models.NewResponseJSON(w, http.StatusCreated, "Professor created successfully", nil)
}

//ASSIGN A COURSE
func (p personHandler) CourseAssignment(w http.ResponseWriter, r *http.Request) {
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		models.NewResponseJSON(w, http.StatusBadRequest, "Invalid json", nil)
		return
	}
	var assigment models.Assignment
	err = json.Unmarshal(req, &assigment)
	if err != nil {
		models.NewResponseJSON(w, http.StatusBadRequest, "Invalid json", nil)
		return
	}

}

//GET ALL THE PEOPLE IN A COURSE FOR A PERIOD
func (p personHandler) GetAllBySectionIDAndDateRange(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	objectID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		models.NewResponseJSON(w, http.StatusBadRequest, "Invalid id", nil)
		return
	}
	startDate, err := time.Parse("2006-01-02", params["startDate"])
	if err != nil {
		models.NewResponseJSON(w, http.StatusBadRequest, "Date format invalid in 'startDate'", nil)
		return
	}
	endDate, err := time.Parse("2006-01-02", params["endDate"])
	if err != nil {
		models.NewResponseJSON(w, http.StatusBadRequest, "Date format invalid in 'endDate'", nil)
		return
	}
	persons, err := p.Persondao.GetAllBySectionID(objectID, startDate, endDate)
	if err != nil {
		models.NewResponseJSON(w, http.StatusInternalServerError, "Server error", nil)
		return
	}
	models.NewResponseJSON(w, http.StatusOK, "OK", persons)
}

func (p personHandler) verifyEmail(email string) bool{
	findPerson := p.Persondao.GetByEmail(email)
	return findPerson == nil
}