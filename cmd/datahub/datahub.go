package main

import (
	"log"
	"net/http"
	"time"

	"github.com/NeowayLabs/datahub"
)

func main() {
	s := &http.Server{
		Addr:           ":8080",
		Handler:        datahub.NewHandler(),
		ReadTimeout:    time.Hour,
		WriteTimeout:   time.Hour,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
