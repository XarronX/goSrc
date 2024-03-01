package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/polyAnalytica/equityConsultancy/src/components/auth"
)

func Establish() {
	router := mux.NewRouter()
	registerRoutes(router)
	fmt.Println("initializing server at :8080...")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}

func registerRoutes(router *mux.Router) {
	router.HandleFunc("/signup", auth.Signup).Methods("POST")
	router.HandleFunc("/login", auth.Login).Methods("POST")
}
