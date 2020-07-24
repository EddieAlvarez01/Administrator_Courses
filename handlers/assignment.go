package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/EddieAlvarez01/administrator_courses/dao"
	"github.com/EddieAlvarez01/administrator_courses/dao/interfaces"
	"github.com/EddieAlvarez01/administrator_courses/models"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
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
	endDate := time.Now().UTC()
	startDate := endDate.AddDate(0, -2, 0)
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

//UPDATE ASSIGNMENT
func (a assignmentHandler) Update(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	assignment, err := a.AssignmentDao.GetById(params["id"])
	if err != nil {
		models.NewResponseJSON(w, http.StatusBadRequest, "Invalid id", nil)
		return
	}
	if assignment == nil {
		models.NewResponseJSON(w, http.StatusNotFound, "The assignment not exist", nil)
		return
	}
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		models.NewResponseJSON(w, http.StatusBadRequest, "Invalid json", nil)
		return
	}
	var editedAssignment models.Assignment
	err = json.Unmarshal(req, &editedAssignment)
	validate := validator.New()
	err = validate.Struct(editedAssignment)
	if err != nil {
		models.NewResponseJSON(w, http.StatusBadRequest, fmt.Sprintf("Errors: %v", err.(validator.ValidationErrors)), nil)
		return
	}

	//ONLY INCLUDE IN THE ARRAY COURSES THAT DON'T EXIST IN THE ASSIGNMENT
	assignment.Courses = a.viewUpdates(assignment.Courses, editedAssignment.Courses)

	err = a.AssignmentDao.Update(assignment)
	if err != nil {
		models.NewResponseJSON(w, http.StatusInternalServerError, "Server error", nil)
		return
	}
	models.NewResponseJSON(w, http.StatusOK, "Assignment updated successfully", assignment)
}

//GET PERIOD ASSIGNMENT
func (a assignmentHandler) GetAssignmentInOnePeriod(w http.ResponseWriter, r *http.Request) {
	var startDate time.Time
	var endDate time.Time
	var err error
	params := mux.Vars(r)
	startDateString := params["startDate"]
	endDateString := params["endDate"]
	if startDateString != "" && endDateString != "" {
		startDate, err = time.Parse("2006-01-02", startDateString)
		if err != nil {
			models.NewResponseJSON(w, http.StatusBadRequest, "Date format invalid in 'startDate'", nil)
			return
		}
		endDate, err = time.Parse("2006-01-02", endDateString)
		if err != nil {
			models.NewResponseJSON(w, http.StatusBadRequest, "Date format invalid in 'endDate'", nil)
			return
		}
	}else{
		endDate = time.Now().UTC()
		startDate = endDate.AddDate(0, -2, 0)
	}
	payload := r.Context().Value("payload").(models.Payload)
	assignment, err := a.AssignmentDao.GetByPersonAndDateRange(payload.ID, startDate, endDate)
	if err != nil {
		models.NewResponseJSON(w, http.StatusInternalServerError, "Server error", nil)
		return
	}
	models.NewResponseJSON(w, http.StatusOK, "OK", assignment)
}

//VALIDATE THAT A SECTION EXIST
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

//VALIDATE THAT THE COURSE NOT EXIST IN THE ASSIGNMENT
func (a assignmentHandler) viewUpdates(assignmentCourses, updateAssignmentCourses []primitive.ObjectID) []primitive.ObjectID {
	var newCourses []primitive.ObjectID
	for _, editedCourseAssignment := range updateAssignmentCourses {
		for i, course := range assignmentCourses {
			if course == editedCourseAssignment {
				newCourses = append(newCourses, course)
				break
			}
			if i == len(assignmentCourses) - 1 {
				newCourses = append(newCourses, editedCourseAssignment)
			}
		}
	}
	return newCourses
}