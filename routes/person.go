package routes

import (
	"net/http"

	"github.com/EddieAlvarez01/administrator_courses/dao/interfaces"
	"github.com/EddieAlvarez01/administrator_courses/handlers"
	"github.com/gorilla/mux"
)

//RegisterRoutesPersons RECORDS THE ROUTES OF PERSONS
func RegisterRoutesPersons(mux *mux.Router, person interfaces.PersonDao) {
	handler := handlers.NewPersonHandler(person)
	mux.HandleFunc("/persons", handler.CreatePerson).Methods(http.MethodPost)
	mux.HandleFunc("/persons/{id}", handler.GetOne).Methods(http.MethodGet)
	mux.HandleFunc("/persons/signin", handler.SignIn).Methods(http.MethodPost)
}
