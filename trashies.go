package main

import (
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"encoding/json"
	"html/template"
)

type Person struct {
	ID        string   `json:"id,omitempty"`
	FirstName string   `json:"firsname,omitempty"`
	LastName  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}

type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

var people []Person

func GetPeople(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(people)
}
func GetPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Person{})
}

func CreatePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	person.ID = params["id"]
	people = append(people, person)
	json.NewEncoder(w).Encode(people)
}

func DeletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range people {
		if item.ID == params["id"] {
			people = append(people[:index], people[index+1:]...)
			break
		}
		json.NewEncoder(w).Encode(people)
	}
}

func launch(w http.ResponseWriter, r *http.Request) {
	baseResponseHandler(w, r, "index.html", "index")
}

func join(w http.ResponseWriter, r *http.Request) {
	baseResponseHandler(w, r, "join.html", "join")
}

func waiting(w http.ResponseWriter, r *http.Request) {
	baseResponseHandler(w, r, "waiting.html", "waiting")
}

func voting(w http.ResponseWriter, r *http.Request) {
	baseResponseHandler(w, r, "voting.html", "voting")
}

func baseResponseHandler(w http.ResponseWriter, r *http.Request, templatePath string, templateName string) {
	baseTemplates := []string{"templates/footer.html", "templates/navbar.html", "templates/header.html", "templates/jsimports.html"}
	baseTemplates = append(baseTemplates, templatePath)
	t, err := template.ParseFiles(baseTemplates...)
	if err != nil {
		print(err)
	}
	t.ExecuteTemplate(w, templateName, nil)
}

func main() {

	cssHandler := http.FileServer(http.Dir("./css/"))
	jsHandler := http.FileServer(http.Dir("./js/"))
	imageHandler := http.FileServer(http.Dir("./images/"))

	http.Handle("/css/", http.StripPrefix("/css/", cssHandler))
	http.Handle("/js/", http.StripPrefix("/js/", jsHandler))
	http.Handle("/images/", http.StripPrefix("/images/", imageHandler))

	people = append(people, Person{ID: "1", FirstName: "John", LastName: "Doe", Address: &Address{City: "City X", State: "State X"}})
	people = append(people, Person{ID: "2", FirstName: "Koko", LastName: "Doe", Address: &Address{City: "City Z", State: "State Y"}})
	people = append(people, Person{ID: "3", FirstName: "Francis", LastName: "Sunday"})

	router := mux.NewRouter()
	router.HandleFunc("/index", launch).Methods("GET")
	router.HandleFunc("/join", join).Methods("GET")
	router.HandleFunc("/waiting", waiting).Methods("GET")
	router.HandleFunc("/voting", voting).Methods("GET")

	router.HandleFunc("/people", GetPeople).Methods("GET")
	router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
	router.HandleFunc("/people/{id}", CreatePerson).Methods("POST")
	router.HandleFunc("/people/{id}", DeletePerson).Methods("POST")
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
