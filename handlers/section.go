package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/EddieAlvarez01/administrator_courses/dao/interfaces"
	"github.com/EddieAlvarez01/administrator_courses/models"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"time"
)

type sectionHandler struct {
	SectionDao interfaces.SectionDao
}

//RETURN A SECTION HANDLER INITIALIZATION
func NewSectionHandler(sectionDao interfaces.SectionDao) sectionHandler{
	return sectionHandler{sectionDao}
}

//CREATE A NEW SECTION
func (s sectionHandler) Create(w http.ResponseWriter, r *http.Request) {
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		models.NewResponseJSON(w, http.StatusBadRequest, "Invalid json", nil)
		return
	}
	var section models.Section
	err = json.Unmarshal(req, &section)
	if err != nil {
		models.NewResponseJSON(w, http.StatusBadRequest, "Invalid json", nil)
		return
	}

	//VALIDATE DATA
	validate := validator.New()
	err = validate.Struct(section)

	//VALIDATE THAT THE INITIAL TIME IS LESS THAN THE FINAL TIME
	if !s.validateSchedule(section.StartHour, section.EndHour) {
		models.NewResponseJSON(w, http.StatusBadRequest, "The start time must be greater than the end time", nil)
		return
	}

	//VALIDATE TWICE SECTIONS REPEAT FOR ONE COURSE
	repeatSection := s.SectionDao.GetByCourseAndSection(section.CourseID, section.Name)
	if repeatSection != nil {
		models.NewResponseJSON(w, http.StatusBadRequest, "You can't repeat sections letters", nil)
		return
	}

	//VALIDATE THE PROFESSOR
	person := s.SectionDao.GetProfessor(section.Professor)
	if person == nil {
		models.NewResponseJSON(w, http.StatusBadRequest, "The professor doesn't exist", nil)
		return
	}
	if !s.validateRole(person.Role) {
		models.NewResponseJSON(w, http.StatusBadRequest, "The person isn't a professor", nil)
		return
	}

	if err != nil {
		models.NewResponseJSON(w, http.StatusBadRequest, fmt.Sprintf("Errors: %v", err.(validator.ValidationErrors)), nil)
		return
	}
	err = s.SectionDao.Create(&section)
	if err != nil {
		models.NewResponseJSON(w, http.StatusInternalServerError, "Error in server", nil)
		return
	}
	models.NewResponseJSON(w, http.StatusCreated, "Section created", nil)
}

//UPDATE A SECTION
func (s sectionHandler) Update(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	sectionToUpdate, err := s.SectionDao.GetById(params["id"])
	if err != nil {
		models.NewResponseJSON(w, http.StatusBadRequest, "Invalid id", nil)
		return
	}
	if sectionToUpdate == nil {
		models.NewResponseJSON(w, http.StatusNotFound, "Section not found", nil)
		return
	}
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		models.NewResponseJSON(w, http.StatusBadRequest, "Invalid json", nil)
		return
	}
	var section models.Section
	err = json.Unmarshal(req, &section)
	if err != nil {
		models.NewResponseJSON(w, http.StatusBadRequest, "Invalid json", nil)
		return
	}
	validate := validator.New()
	err = validate.Struct(section)
	if err != nil {
		models.NewResponseJSON(w, http.StatusBadRequest, fmt.Sprintf("Errors: %v", err.(validator.ValidationErrors)), nil)
		return
	}
	if sectionToUpdate.Professor != section.Professor {

		//VALIDATE EXIST PROFESSOR AND ROLE PROFESSOR
		person := s.SectionDao.GetProfessor(section.Professor)
		if person == nil {
			models.NewResponseJSON(w, http.StatusNotFound, "The professor doesn't exist", nil)
			return
		}
		if !s.validateRole(person.Role) {
			models.NewResponseJSON(w, http.StatusBadRequest, "The person isn't a professor", nil)
			return
		}
		sectionToUpdate.Professor = section.Professor

	}
	var flag = false
	if sectionToUpdate.CourseID != section.CourseID {
		sectionToUpdate.CourseID = section.CourseID
		flag = true
	}
	if sectionToUpdate.Name != section.Name {
		sectionToUpdate.Name = section.Name
		flag = true
	}
	sectionToUpdate.StartHour = section.StartHour
	sectionToUpdate.EndHour = section.EndHour

	//VALIDATE SCHEDULE
	if !s.validateSchedule(sectionToUpdate.StartHour, sectionToUpdate.EndHour) {
		models.NewResponseJSON(w, http.StatusBadRequest, "The start time must be greater than the end time", nil)
		return
	}

	//VALIDATE DUPLICATE SECTIONS
	if exist := s.SectionDao.GetByCourseAndSection(sectionToUpdate.CourseID, sectionToUpdate.Name); flag && exist != nil {
		models.NewResponseJSON(w, http.StatusBadRequest, "The course cannot have duplicate sections", nil)
		return
	}

	err = s.SectionDao.Update(sectionToUpdate)
	if err != nil {
		models.NewResponseJSON(w, http.StatusInternalServerError, "Error in server", nil)
		return
	}
	models.NewResponseJSON(w, http.StatusOK, "OK", sectionToUpdate)
}

//VALIDATE COURSE SCHEDULE
func (s sectionHandler) validateSchedule(initial time.Time, end time.Time) bool {
	initialMinutes := initial.Hour() * 60 + initial.Minute()
	endMinutes := end.Hour() * 60 + end.Minute()
	return initialMinutes < endMinutes
}

//VALIDATE PERSON ROLE
func (s sectionHandler) validateRole(roles []string) bool {
	for _, role := range roles {
		if role == "PROFESSOR" {
			return true
		}
	}
	return false
}