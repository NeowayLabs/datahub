package datahub

import (
	"log"
	"net/http"
	"os"
)

type handler struct {
	datadir  string
	datafile string
	log      *log.Logger
}

func (d *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/datahub/upload":
		{
			d.upload(w, req)
		}
	default:
		{
			d.log.Printf("error: path %q not found", req.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func NewHandler() http.Handler {
	return &handler{
		datadir: "./.datafiles",
		log:     log.New(os.Stdout, "datahub.server", log.Lshortfile|log.Lmicroseconds),
	}
}

func (d *handler) upload(w http.ResponseWriter, req *http.Request) {
}
