package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

type UserDetail struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func testFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Content-Type", "application/json")
	var newUser UserDetail
	var response Response
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		fmt.Println(err)
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}
	fmt.Printf("%+v", newUser)
	response.Message = "Successfully submitted"
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/test", testFunc)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
