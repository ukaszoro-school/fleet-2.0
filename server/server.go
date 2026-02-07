package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"context"
	"log"
	"time"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	uri := "mongodb://localhost:27017"

	s, client := storageNew(&ctx, uri)

	fmt.Println("Using collection:", s.userCollection.Name())
	fmt.Println("Using collection:", s.routeCollection.Name())
	fmt.Println("Using collection:", s.stopCollection.Name())

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()
	setRoutes(s)

	err := http.ListenAndServe(":8080", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
