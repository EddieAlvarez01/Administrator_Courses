package routes

import (
	"github.com/EddieAlvarez01/administrator_courses/dao"
	"github.com/EddieAlvarez01/administrator_courses/handlers"
	"github.com/EddieAlvarez01/administrator_courses/middlewares"
	"github.com/gorilla/mux"
	"net/http"
)

//RegisterRoutesCourses REGISTER ALL ROUTES FOR COURSES
func RegisterRoutesCourses(mux *mux.Router, courseImpl dao.CourseImpl) {
	handler := handlers.NewCourseHandler(courseImpl)
	mux.Handle("/courses", middlewares.Authenticate(middlewares.RoleAdministrator(http.HandlerFunc(handler.Create)))).Methods(http.MethodPost)
	mux.Handle("/courses/{id}", middlewares.Authenticate(middlewares.RoleAdministrator(http.HandlerFunc(handler.Update)))).Methods(http.MethodPut)
}
