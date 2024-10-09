package main

import (
	"log/slog"
	"net/http"

	"github.com/couryrr/go-rest-persistance/internal/person"
)

func main(){
	personHandler := person.NewPersonHandler()
	router := http.NewServeMux()

	router.HandleFunc(personHandler.HandleAddPerson())

	slog.Error("server stopped","reason",http.ListenAndServe(":8080", router))
}

