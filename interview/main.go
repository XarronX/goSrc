package main

import (
	"interview/middleware"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/allow", middleware.Allow)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
