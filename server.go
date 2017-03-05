package datahub

import (
	"bufio"
	"encoding/json"
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
			uploadFileNames := []string{
				"trainingset.csv",
				"testset.challenge.csv",
				"testset.prediction.csv",
				"testset.result.csv",
			}

			success := false
			for _, filename := range uploadFileNames {
				res := d.receiveUpload(req, filename)
				if res {
					success = res
				}
			}
			if !success {
				d.failrequest(
					w,
					"no dataset received, expected one of these: %q",
					uploadFileNames,
				)
				return
			}
			w.WriteHeader(http.StatusOK)
		}
	case "/datahub/score":
		{
			d.scoreCheck(w, req)
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

func (d *handler) scoreCheck(w http.ResponseWriter, req *http.Request) float32 {
	const datadir string = "./.repo"
	predictionfile, err := os.Open(datadir + "/testset.prediction.csv")
	if err != nil {
		d.log.Printf("error: %q reading file", err)
	}
	defer predictionfile.Close()

	resultfile, err := os.Open(datadir + "/testset.result.csv")
	if err != nil {
		d.log.Printf("error: %q reading file", err)
	}
	defer resultfile.Close()

	scanpredictionfile := bufio.NewScanner(predictionfile)
	scanresultfile := bufio.NewScanner(resultfile)

	totallines := 0
	ok := 0
	for scanresultfile.Scan() {
		totallines = totallines + 1
		if scanpredictionfile.Scan() {
			if scanresultfile.Text() == scanpredictionfile.Text() {
				ok++
			}
		}
	}
	score := float32(int((ok*100/totallines)*100)) / 100
	response := make(map[string]interface{})
	response["score"] = score

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(response)
	return score
}
