package main

import (
	"fmt"
	"io"
	"net/http"
	"text/template"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func isHTMX(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

func getRoot() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		tmpl.Execute(w, nil)
	})
}

func getHello() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		io.WriteString(w, "Hello, HTTP!\n")
	})
}

func stopPageHandler(s *Storage, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if isHTMX(r) {

				var stop []Stop

				stop, _ = s.getAllStops()

				tmpl.ExecuteTemplate(w, "stops.html", stop)
			} else {

				tmpl := template.Must(template.ParseFiles("templates/index.html"))
				tmpl.Execute(w, nil)
			}
		case http.MethodPost:
			location := r.FormValue("location")
			if location == "" {
				http.Error(w, "location is empty", http.StatusBadRequest)
				return
			}

			stop := Stop{
				ID:       primitive.NewObjectID(),
				Location: location,
			}
			_, err := s.createStop(&stop)
			if err != nil {
				fmt.Print(err)
			}
			if err := tmpl.ExecuteTemplate(w, "stop-row", stop); err != nil {
				http.Error(w, "Template error", http.StatusInternalServerError)
				return
			}
		}
	}
}

func stopDeleteHandler(s *Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idstr := r.PathValue("id")
		fmt.Print(r.URL)
		id, err := primitive.ObjectIDFromHex(idstr)
		if err != nil {
			fmt.Print("Failed to delete object: ", err)

		}
		s.deleteStopByID(id)
	}
}
