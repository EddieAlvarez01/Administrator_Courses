package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/EddieAlvarez01/administrator_courses/dao"

	"github.com/EddieAlvarez01/administrator_courses/routes"
	"github.com/gorilla/mux"
)

func main() {

	//CONFIG ROUTES
	r := mux.NewRouter()
	router := r.PathPrefix("/api").Subrouter()
	routes.RegisterRoutesPersons(router, dao.PersonImpl{})

	//SERVER
	fmt.Println("Server on port 7000")
	log.Fatal(http.ListenAndServe(":7000", router))

}
