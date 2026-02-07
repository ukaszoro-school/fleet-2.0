package main

import (
	"net/http"
	"strings"
	"text/template"
)

func FileServerFilter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Disable dir listing
		if strings.HasSuffix(r.URL.Path, "/") && r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func ServeFile(name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, name)
	})
}

func setRoutes(s *Storage) {
	tmpl := template.Must(template.ParseGlob("templates/*.html"))
	fs := http.FileServer(http.Dir("../content/"))
	http.Handle("GET /content/", http.StripPrefix("/content", FileServerFilter(fs)))

	http.Handle("/", getRoot())
	http.Handle("GET /hello", getHello())
	http.Handle("/stops", stopPageHandler(s, tmpl))
	http.Handle("/routes", routePageHandler(s, tmpl))
	http.Handle("GET /routes/new-row", routeRowPageHandler(s, tmpl))
	http.Handle("DELETE /stops/{id}", stopDeleteHandler(s))
	http.Handle("DELETE /routes/{id}", routeDeleteHandler(s))
	http.Handle("/home", getHello())

}
