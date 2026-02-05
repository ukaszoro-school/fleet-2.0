package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"context"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
	"time"
)

type driver struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Age     uint   `json:"age"`
}

var drivers = []driver{
	{ID: "1", Name: "Harold", Surname: "Mason", Age: 52},
	{ID: "2", Name: "Leticia", Surname: "Alvarez", Age: 34},
	{ID: "3", Name: "Bojan", Surname: "Petrovic", Age: 45},
}

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

func setRoutes() {
	fs := http.FileServer(http.Dir("../content/"))
	http.Handle("GET /content/", http.StripPrefix("/content", FileServerFilter(fs)))

	http.Handle("/", getRoot())
	http.Handle("GET /hello", getHello())

}

func getRoot() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("got / request\n")
		io.WriteString(w, "This is my website!\n")
	})
}

func getHello() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("got /hello request\n")
		io.WriteString(w, "Hello, HTTP!\n")
	})
}

func main() {
	setRoutes()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := "mongodb://localhost:27017"

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Could not connect to MongoDB:", err)
	}

	fmt.Println("Connected to MongoDB!")

	db := client.Database("testdb")
	collection := db.Collection("users")

	fmt.Println("Using database:", db.Name())
	fmt.Println("Using collection:", collection.Name())

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	err = http.ListenAndServe(":8080", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
