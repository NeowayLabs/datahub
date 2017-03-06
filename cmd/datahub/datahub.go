package main

import (
	"log"

	"github.com/NeowayLabs/datahub"
)

func main() {
	server := datahub.NewServer()

	log.Fatal(Server.ListenAndServe(":8080"))
}
