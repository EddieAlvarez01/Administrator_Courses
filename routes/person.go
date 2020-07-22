package routes

import (
	"net/http"

	"github.com/EddieAlvarez01/administrator_courses/dao/interfaces"
	"github.com/EddieAlvarez01/administrator_courses/handlers"
	"github.com/EddieAlvarez01/administrator_courses/middlewares"
	"github.com/gorilla/mux"
)

//RegisterRoutesPersons RECORDS THE ROUTES OF PERSONS
func RegisterRoutesPersons(mux *mux.Router, person interfaces.PersonDao) {
	handler := handlers.NewPersonHandler(person)
	mux.HandleFunc("/persons", handler.CreatePerson).Methods(http.MethodPost)
	mux.HandleFunc("/persons/{id}", handler.GetOne).Methods(http.MethodGet)
	mux.HandleFunc("/persons/signin", handler.SignIn).Methods(http.MethodPost)
	mux.Handle("/persons/update", middlewares.Authenticate(http.HandlerFunc(handler.Update))).Methods(http.MethodPut)
	mux.Handle("/persons/new-professor", middlewares.Authenticate(middlewares.PersonRole(http.HandlerFunc(handler.CreateProfessor), 0))).Methods(http.MethodPost)
}
