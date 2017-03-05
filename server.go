package datahub

import (
	"io"
	"log"
	"net/http"
	"os"
)

type handler struct {
	datadir string
	log     *log.Logger
}

func (d *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// NOT VALIDATING THE HTTP METHODS :-)
	switch req.URL.Path {
	case "/datahub/upload":
		{
			trainingset := d.receiveUpload(req, "trainingset")
			testset := d.receiveUpload(req, "testset")
			testsetres := d.receiveUpload(req, "testsetres")
			if !trainingset && !testset && !testsetres {
				d.failrequest(w, "no dataset received, expected one of these: testset, testsetres, trainingset")
				return
			}
			w.WriteHeader(http.StatusOK)
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
		datadir: datadir,
		log:     log,
	}
}

func (d *handler) failrequest(
	w http.ResponseWriter,
	fmt string,
	args ...interface{},
) {
	d.log.Printf(fmt, args...)
	w.WriteHeader(http.StatusInternalServerError)
}

func (d *handler) receiveUpload(
	req *http.Request,
	filename string,
) bool {
	uploadedfile, _, err := req.FormFile(filename)
	if err != nil {
		d.log.Printf("%q parsing form", err)
		return false
	}
	defer uploadedfile.Close()

	filepath := d.datadir + "/" + filename
	d.log.Printf("creating file %q", filepath)
	file, err := os.Create(filepath)
	if err != nil {
		d.log.Printf("error: %q opening file %q", err, filepath)
		return false
	}
	d.log.Printf("created file with success, copying contents")
	_, err = io.Copy(file, uploadedfile)
	if err != nil {
		d.log.Printf("error: %q copying file", err)
		return false
	}
	d.log.Printf("finished copying from form %q with success", filename)
	return true
}
