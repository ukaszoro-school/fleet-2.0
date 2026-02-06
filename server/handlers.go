package main

import (
	"fmt"
	"io"
	"net/http"
)

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
