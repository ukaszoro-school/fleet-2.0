package main

import (
	"io"
	"net/http"
	"text/template"
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
		if isHTMX(r) {

			var stop []Stop

			stop, _ = s.getAllStops()

			tmpl.ExecuteTemplate(w, "stops.html", stop)
		} else {

			tmpl := template.Must(template.ParseFiles("templates/index.html"))
			tmpl.Execute(w, nil)
		}

	}
}
