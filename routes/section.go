package routes

import (
	"github.com/EddieAlvarez01/administrator_courses/dao"
	"github.com/EddieAlvarez01/administrator_courses/handlers"
	"github.com/EddieAlvarez01/administrator_courses/middlewares"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterRoutesSection(muxer *mux.Router, sectionImpl dao.SectionImpl){
	handler := handlers.NewSectionHandler(sectionImpl)
	muxer.Handle("/sections", middlewares.Authenticate(middlewares.PersonRole(http.HandlerFunc(handler.Create), 0))).Methods(http.MethodPost)
	muxer.Handle("/sections/{id}", middlewares.Authenticate(middlewares.PersonRole(http.HandlerFunc(handler.Update), 0))).Methods(http.MethodPut)
}
