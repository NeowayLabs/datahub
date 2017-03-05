package datahub

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
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
				"code.r",
				"trainingset.csv",
				"testset.challenge.csv",
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
	case "/datahub/execr":
		{
			d.execR(w, req)
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

func (d *handler) execR(
	w http.ResponseWriter,
	req *http.Request,
) {
	// TODO: Still not getting stderr
	cwd, err := os.Getwd()
	if err != nil {
		d.failrequest(w, "getwd: unexpected error %q", err)
		return
	}
	err = os.Chdir(cwd + "/" + d.datadir)
	if err != nil {
		d.failrequest(w, "chdir: unexpected error %q", err)
		return
	}
	defer os.Chdir(cwd)

	cmd := exec.Command("R", "-f", "./code.r")
	res, err := cmd.CombinedOutput()

	if err != nil {
		d.failrequest(w, "exec R: unexpected error %q", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	if err != nil {
		d.log.Printf("unexpected error %q sending response")
		return
	}
	d.log.Printf("executed R code with success")
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
