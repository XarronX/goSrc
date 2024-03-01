package server

import (
	"fmt"
	"net/http"

	"github.com/ashutosh/goethTrieAPIs/api"
	"github.com/gorilla/mux"
)

func Server() {
	r := mux.NewRouter()

	r.HandleFunc("/setvalue", api.SetValue).Methods("POST")
	r.HandleFunc("/getvalue", api.GetValue).Methods("POST")
	r.HandleFunc("/gethash", api.GetHash).Methods("GET")
	r.HandleFunc("/verifyhash", api.MerkleVerification).Methods("POST")
	r.HandleFunc("/securetrie", api.SecureTree).Methods("GET")
	r.HandleFunc("/securetrie/getvalue", api.GetSecureTreeVal).Methods("POST")

	fmt.Println("Starting server :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}
