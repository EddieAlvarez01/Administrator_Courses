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
)

type courseHandler struct {
	CourseDao interfaces.CourseDao
}

//NewCourseHandler RETURN INSTANCE COURSE HANDLER
func NewCourseHandler(courseDao interfaces.CourseDao) courseHandler{
	return courseHandler{courseDao}
}

//Create (A NEW COURSE)
func (c courseHandler) Create(w http.ResponseWriter, r *http.Request) {
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		models.NewResponseJSON(w, http.StatusBadRequest, "Invalid json", nil)
		return
	}
	var course models.Course
	json.Unmarshal(req, &course)
	validate := validator.New()
	if err = validate.Struct(course); err != nil {
		models.NewResponseJSON(w, http.StatusBadRequest, fmt.Sprintf("Invalid data: %v", err.(validator.ValidationErrors)), nil)
		return
	}
	err = c.CourseDao.Create(&course)
	if err != nil {
		models.NewResponseJSON(w, http.StatusInternalServerError, "Error inserting the course", nil)
		return
	}
	models.NewResponseJSON(w, http.StatusCreated, "Course created successfully", course)
}

//UPDATE A COURSE
func (c courseHandler) Update(w http.ResponseWriter, r *http.Request) {
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		models.NewResponseJSON(w, http.StatusBadRequest, "Invalid json", nil)
		return
	}
	var updateCourse models.Course
	json.Unmarshal(req, &updateCourse)
	validate := validator.New()
	if err = validate.Struct(updateCourse); err != nil {
		models.NewResponseJSON(w, http.StatusBadRequest, fmt.Sprintf("Invalid data: %v", err.(validator.ValidationErrors)), nil)
		return
	}
	params := mux.Vars(r)
	course, err := c.CourseDao.GetById(params["id"])
	if err != nil {
		models.NewResponseJSON(w, http.StatusBadRequest, "Invalid id", nil)
		return
	}
	if course == nil {
		models.NewResponseJSON(w, http.StatusBadRequest, "Course not found", nil)
		return
	}
	course.Name = updateCourse.Name
	course.Credits = updateCourse.Credits
	course.School = updateCourse.School
	err = c.CourseDao.Update(course)
	if err != nil {
		models.NewResponseJSON(w, http.StatusInternalServerError, "Error on update course", nil)
		return
	}
	models.NewResponseJSON(w, http.StatusOK, "OK", course)
}

