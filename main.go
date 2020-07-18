package main

import (
	"fmt"
	"github.com/EddieAlvarez01/administrator_courses/dao"
	"github.com/EddieAlvarez01/administrator_courses/routes"
	"github.com/gorilla/mux"
	"log"
	"net/http"

	"github.com/EddieAlvarez01/administrator_courses/authorization"
)

func main() {

	//KEYS JWT
	err := authorization.LoadFiles("certificates/app.rsa", "certificates/app.rsa.pub")
	if err != nil {
		log.Fatal(err)
	}

	//CONFIG ROUTES
	r := mux.NewRouter()
	router := r.PathPrefix("/api").Subrouter()
	routes.RegisterRoutesPersons(router, dao.PersonImpl{})
	routes.RegisterRoutesCourses(router, dao.CourseImpl{})

	//SERVER
	fmt.Println("Server on port 7000")
	log.Fatal(http.ListenAndServe(":7000", router))
}
