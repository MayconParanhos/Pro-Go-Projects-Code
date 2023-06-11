package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Rsvp struct {
	Name, Email, Phone string
	WillAttend         bool
}

type formData struct {
	*Rsvp
	Errors []string
}

type response struct {
	Name string
	Resp *[]*Rsvp
}

var formEmpty = formData{Rsvp: &Rsvp{}, Errors: []string{}}

var responses = make([]*Rsvp, 0, 10)
var templates = make(map[string]*template.Template, 3)

func loadTemplates() {
	var templateNames = [5]string{"welcome", "form", "thanks", "sorry", "list"}
	for index, name := range templateNames {
		templt, err := template.ParseFiles("layout.html", name+".html")
		if err == nil {
			templates[name] = templt
			fmt.Println("Loaded templates:", index, name)
		} else {
			panic(err)
		}
	}
}

func welcomeHandler(writer http.ResponseWriter, request *http.Request) {
	templates["welcome"].Execute(writer, nil)
}

func listHandler(writer http.ResponseWriter, request *http.Request) {
	templates["list"].Execute(writer, responses)
}

func formHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		templates["form"].Execute(writer, &formEmpty)
	} else if request.Method == http.MethodPost {
		request.ParseForm()
		responseData := Rsvp{
			Name:       request.Form["name"][0],
			Email:      request.Form["email"][0],
			Phone:      request.Form["phone"][0],
			WillAttend: request.Form["willattend"][0] == "true",
		}

		var errors []string
		if responseData.Name == "" {
			errors = append(errors, "Por favor, coloque o seu nome.")
		}
		if responseData.Email == "" {
			errors = append(errors, "Por favor, coloque o seu endereço de email.")
		}
		if responseData.Phone == "" {
			errors = append(errors, "Por favor, coloque o seu número de telefone.")
		}
		if len(errors) > 0 {
			templates["form"].Execute(writer, formData{Rsvp: &responseData, Errors: errors})
		} else {
			if responseData.WillAttend {
				responses = append(responses, &responseData)
				templates["thanks"].Execute(writer, responseData.Name)
			} else {
				responses = append(responses, &responseData)
				templates["sorry"].Execute(writer, response{Name: responseData.Name, Resp: &responses})
			}
		}
	}
}

func main() {
	go loadTemplates()
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/form", formHandler)

	log.Fatalln(http.ListenAndServe(":8080", nil))
}
