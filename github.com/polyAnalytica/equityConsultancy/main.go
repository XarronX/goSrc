package main

import (
	"log"

	"github.com/polyAnalytica/equityConsultancy/src/db"
	"github.com/polyAnalytica/equityConsultancy/src/server"
)

func main() {
	_, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	server.Establish()
}
