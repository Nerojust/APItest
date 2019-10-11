/*package main

import (
	"fmt"
	"math"
	"strconv"
)

func main() {
	//person1 := Person{
	//	LastName:  "Alberta",
	//	FirstName: "Gladys",
	//	gender:    "female",
	//	Age:       20,
	//	Address: Address{
	//		State: "Lagos",
	//		City:  "Ipaja",
	//	},
	//}
	//person1.getMarried("Travis")
	//fmt.Println(person1.greet())

}
func (circle Circle) getCircleArea() float64 {
	return math.Pi * circle.x * circle.y * circle.radius
}

func (rect Rectangle) getRectangleArea() float64 {
	return rect.height * rect.width
}

func getArea(shape ShapeInterface) float64 {
	return shape.area()
}

type ShapeInterface interface {
	area() float64
}
type Circle struct {
	x, y, radius float64
}
type Rectangle struct {
	height, width float64
}

func (person *Person) getMarried(lastName string) {
	if person.gender == "male" {
		fmt.Println("Its a male")
		return
	} else {
		person.LastName = lastName
	}
}
func (p Person) greet() string {
	return "Hello my name is " + p.FirstName + " " + p.LastName + " and i am " + strconv.Itoa(p.Age) + " years old" + " I live at " + p.Address.City + ", " + p.Address.State
}

type Person struct {
	LastName, FirstName, gender string
	Age                         int
	Address                     Address
}
type Address struct {
	State, City string
}
*/


package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type event struct {
	ID          string `json:"ID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

type allEvents []event

var events = allEvents{
	{
		ID:          "1",
		Title:       "Introduction to Golang",
		Description: "Come join us for a chance to learn how golang works and get to eventually try it out",
	},
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent event
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}

	json.Unmarshal(reqBody, &newEvent)
	events = append(events, newEvent)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newEvent)
}

func getOneEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	for _, singleEvent := range events {
		if singleEvent.ID == eventID {
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}

func getAllEvents(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(events)
}

func updateEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]
	var updatedEvent event

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}
	json.Unmarshal(reqBody, &updatedEvent)

	for i, singleEvent := range events {
		if singleEvent.ID == eventID {
			singleEvent.Title = updatedEvent.Title
			singleEvent.Description = updatedEvent.Description
			events = append(events[:i], singleEvent)
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}

func deleteEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	for i, singleEvent := range events {
		if singleEvent.ID == eventID {
			events = append(events[:i], events[i+1:]...)
			fmt.Fprintf(w, "The event with ID %v has been deleted successfully", eventID)
		}
	}
}

func main() {
	//initEvents()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/event", createEvent).Methods("POST")
	router.HandleFunc("/events", getAllEvents).Methods("GET")
	router.HandleFunc("/events/{id}", getOneEvent).Methods("GET")
	router.HandleFunc("/events/{id}", updateEvent).Methods("PATCH")
	router.HandleFunc("/events/{id}", deleteEvent).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}