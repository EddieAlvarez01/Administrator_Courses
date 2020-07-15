package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/EddieAlvarez01/administrator_courses/models"

	"github.com/EddieAlvarez01/administrator_courses/dao/interfaces"

	"github.com/go-playground/validator/v10"
)

type personHandler struct {
	Persondao interfaces.PersonDao
}

//NewPersonHandler return a new instance of NewPersonHandler
func NewPersonHandler(person interfaces.PersonDao) personHandler {
	return personHandler{person}
}

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
	err = p.Persondao.Create(&person)
	if err != nil {
		models.NewResponseJSON(w, http.StatusInternalServerError, fmt.Sprintf("Server error: %s", err.Error()), nil)
		return
	}
	models.NewResponseJSON(w, http.StatusCreated, "Person created successfully", nil)
}
