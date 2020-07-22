package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/EddieAlvarez01/administrator_courses/dao"
	"github.com/EddieAlvarez01/administrator_courses/dao/interfaces"
	"github.com/EddieAlvarez01/administrator_courses/models"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"net/http"
	"time"
)

type assignmentHandler struct {
	AssignmentDao interfaces.AssignmentDao
}

//RETURN A NEW INSTANCE OF ASSIGNMENT HANDLER
func NewAssignmentHandler (assignmentDao interfaces.AssignmentDao) assignmentHandler{
	return assignmentHandler{assignmentDao}
}

//CREATE A NEW ASSIGMENT
func (a assignmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		models.NewResponseJSON(w, http.StatusBadRequest, "Invalid json", nil)
		return
	}
	var assignment models.Assignment
	err = json.Unmarshal(req, &assignment)
	if err != nil {
		models.NewResponseJSON(w, http.StatusBadRequest, "Invalid json", nil)
		return
	}
	validate := validator.New()
	err = validate.Struct(assignment)
	if err != nil {
		models.NewResponseJSON(w, http.StatusBadRequest, fmt.Sprintf("Errors: %v", err.(validator.ValidationErrors)), nil)
		return
	}
	payload := r.Context().Value("payload").(models.Payload)

	//VALIDATE THAT YOU DON'T CREATE TWO ASSIGNMENTS IN ONE PERIOD
	startDate := time.Now().UTC()
	endDate := startDate.AddDate(0, -2, 0)
	findAssignment, err := a.AssignmentDao.GetByPersonAndDateRange(payload.ID, startDate, endDate)
	if err != nil {
		models.NewResponseJSON(w, http.StatusBadRequest, "Invalid person id", nil)
		return
	}
	if findAssignment != nil {
		models.NewResponseJSON(w, http.StatusBadRequest, "You can't have a two assignments in one period", nil)
		return
	}

	//VALIDATE THAT THE COURSES EXIST
	if !a.validateSections(assignment.Courses) {
		models.NewResponseJSON(w, http.StatusNotFound, "A course was not found", nil)
		return
	}

	assignment.Date = time.Now().UTC()
	assignment.PersonID, _ = primitive.ObjectIDFromHex(payload.ID)
	err = a.AssignmentDao.Create(&assignment)
	if err != nil {
		models.NewResponseJSON(w, http.StatusInternalServerError, "Server error", nil)
		return
	}
	models.NewResponseJSON(w, http.StatusCreated, "OK", assignment)
}

func (a assignmentHandler) validateSections(courses []primitive.ObjectID) bool {
	var sectionDao interfaces.SectionDao = dao.SectionImpl{}
	for _, course := range courses {
		section, err := sectionDao.GetById(course.Hex())
		if err != nil || section == nil {
			return false
		}
	}
	return true
}