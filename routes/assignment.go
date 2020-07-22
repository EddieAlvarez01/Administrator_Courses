package routes

import (
	"github.com/EddieAlvarez01/administrator_courses/dao"
	"github.com/EddieAlvarez01/administrator_courses/handlers"
	"github.com/EddieAlvarez01/administrator_courses/middlewares"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterRoutesAssignment(mux *mux.Router, assignmentImpl dao.AssignmentImpl) {
	handler := handlers.NewAssignmentHandler(assignmentImpl)
	mux.Handle("/assignments", middlewares.Authenticate(middlewares.PersonRole(http.HandlerFunc(handler.Create), 1))).Methods(http.MethodPost)
}
