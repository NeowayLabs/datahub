package datahub

import (
	"io"
	"log"
	"net/http"
	"os"
)

type handler struct {
	datafile string
	log      *log.Logger
}

func (d *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// NOT VALIDATING THE HTTP METHODS :-)
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
	const datadir string = "./.repo"
	log := log.New(os.Stdout, "datahub.server", log.Lshortfile|log.Lmicroseconds)
	err := os.MkdirAll(datadir, 0755)
	if err != nil {
		log.Fatalf("error %q creating data dir %q", err, datadir)
	}
	return &handler{
		datafile: datadir + "/uploaded_data",
		log:      log,
	}
}

func (d *handler) upload(w http.ResponseWriter, req *http.Request) {
	d.log.Printf("creating file %q", d.datafile)
	file, err := os.Create(d.datafile)
	if err != nil {
		d.log.Printf("error: %q opening file %q", err, d.datafile)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	d.log.Printf("created file with success, copying contents")
	_, err = io.Copy(file, req.Body)
	if err != nil {
		d.log.Printf("error: %q copying file", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	d.log.Printf("finished copying")
	w.WriteHeader(http.StatusOK)
}
