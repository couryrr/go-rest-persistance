package person

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Person struct {
	name string
	age int
	height float32
}

type PersonHandler struct{}

type handlerFunc func(w http.ResponseWriter, r *http.Request) 


func NewPersonHandler() *PersonHandler {
	return &PersonHandler{}
}
func (h *PersonHandler) HandleAddPerson()(string, handlerFunc){
	return "POST /person" , func(w http.ResponseWriter, r *http.Request) {
		var person Person
		err := json.NewDecoder(r.Body).Decode(&person)
		if err != nil {
			slog.Error("unable to decode person", "err", err.Error())
		}
		slog.Info("create person", "person", person)
	}
}
func handleAddPepole(people []Person){}
func handleGetPeople(){}
func handleGetPersonById(id int){}
func handleGetPersonByName(name string){}


