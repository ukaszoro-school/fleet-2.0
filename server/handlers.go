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

		io.WriteString(w, "Welcome user to fleet manager home!\n")
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

func routePageHandler(s *Storage, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if isHTMX(r) {

				var stops []Stop
				var routes []Route

				stops, _ = s.getAllStops()
				routes, _ = s.getAllRoutes()

				tmpl.ExecuteTemplate(w, "routes.html", struct {
					Stops  []Stop
					Routes []Route
				}{Stops: stops, Routes: routes})

			} else {

				tmpl := template.Must(template.ParseFiles("templates/index.html"))
				tmpl.Execute(w, nil)
			}
		case http.MethodPost:
			routeName := r.FormValue("route_name")
			stops := r.Form["stop_id[]"]
			times := r.Form["time[]"]

			timesMap := make(map[string]string)

			for i := range stops {
				if stops[i] == "" {
					continue
				}
				timesMap[stops[i]] = times[i]
			}

			route := Route{
				ID:    primitive.NewObjectID(),
				Name:  routeName,
				Times: timesMap,
			}
			fmt.Print(route)
			_, err := s.createRoute(&route)
			if err != nil {
				fmt.Print(err)
			}
			if err := tmpl.ExecuteTemplate(w, "route-read-row", route); err != nil {
				http.Error(w, "Template error", http.StatusInternalServerError)
				return
			}

		}
	}
}

func routeRowPageHandler(s *Storage, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var stop []Stop

		stop, _ = s.getAllStops()
		tmpl.ExecuteTemplate(w, "route_form_row.html", struct {
			Stops []Stop
		}{Stops: stop})
	}
}

func routeDeleteHandler(s *Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idstr := r.PathValue("id")
		fmt.Print(r.URL)
		id, err := primitive.ObjectIDFromHex(idstr)
		if err != nil {
			fmt.Print("Failed to delete object: ", err)

		}
		s.deleteRouteByID(id)
	}
}

func linePageHandler(s *Storage, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if isHTMX(r) {

				var stops []Stop
				var routes []Route

				stops, _ = s.getAllStops()
				routes, _ = s.getAllRoutes()

				stopMap := make(map[string]Stop)
				for _, s := range stops {
					stopMap[s.ID.Hex()] = s
				}

				tmpl.ExecuteTemplate(w, "lines.html", struct {
					StopMap map[string]Stop
					Routes  []Route
				}{StopMap: stopMap, Routes: routes})

			} else {

				tmpl := template.Must(template.ParseFiles("templates/index.html"))
				tmpl.Execute(w, nil)
			}

		}
	}
}
