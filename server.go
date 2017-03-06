package datahub

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/NeowayLabs/datahub/company"
	"github.com/julienschmidt/httprouter"
)

// Server ...
type Server struct {
	router  *httprouter.Router
	datadir string
	log     *log.Logger
	company *company.Company
}

// NewServer ...
func NewServer() *Server {
	const datadir string = "./.repo"
	log := log.New(os.Stdout, "datahub.server", log.Lshortfile|log.Lmicroseconds)
	err := os.MkdirAll(datadir, 0755)
	if err != nil {
		log.Fatalf("error %q creating data dir %q", err, datadir)
	}

	router := httprouter.New()
	company := company.NewCompany()

	d := &Server{
		router:  router,
		datadir: datadir,
		log:     log,
		company: company,
	}

	router.GET("/api/companies/jobs", d.companiesGetJobs)

	router.POST("/api/companies/jobs", d.companiesCreateJob)
	router.GET("/api/companies/jobs/:id", d.companiesGetJob)
	router.POST("/api/companies/job/:id/upload", d.companiesUploadJob)
	router.POST("/api/companies/jobs/:id/start", d.companiesStartJob)

	router.GET("/api/scientists", d.scientistsList)

	router.GET("/api/scientists/:id/jobs", d.scientistsGetJobs)
	router.POST("/api/scientists/:id/jobs/:job/apply", d.scientistsApplyJob)
	router.GET("/api/scientists/:id/jobs/:job/workspace", d.scientistsGetWorkspace)
	router.POST("/api/scientists/:id/jobs/:job/upload", d.companiesUploadCode)

	router.POST("/api/execR", d.execR)

	return d
}

func (d *Server) companiesUploadJob(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	uploadFileNames := []string{
		"code.r",
		"trainingset.csv",
		"testset.challenge.csv",
		"testset.result.csv",
	}

	jobID := params.ByName("id")

	success := false
	for _, filename := range uploadFileNames {
		res := d.receiveUpload(req, filename, jobID)
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

// Jobs ...
type Jobs struct {
	Pending []company.Job `json:"pending"`
	Doing   []company.Job `json:"doing"`
	Done    []company.Job `json:"done"`
}

func (d *Server) companiesGetJobs(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	pending := d.company.GetJobsByStatus("pending")
	doing := d.company.GetJobsByStatus("doing")
	done := d.company.GetJobsByStatus("done")

	jobs := &Jobs{
		Pending: pending,
		Doing:   doing,
		Done:    done,
	}

	bytes, err := json.Marshal(jobs)
	if err != nil {
		d.log.Printf("marshal: error %q", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(bytes); err != nil {
		d.log.Printf("write: error %q", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (d *Server) companiesCreateJob(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	//TODO
	w.WriteHeader(http.StatusNotImplemented)
}

func (d *Server) companiesGetJob(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	//TODO
	w.WriteHeader(http.StatusNotImplemented)
}

func (d *Server) companiesStartJob(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	//TODO
	w.WriteHeader(http.StatusNotImplemented)
}

func (d *Server) scientistsList(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	//TODO
	w.WriteHeader(http.StatusNotImplemented)
}

func (d *Server) scientistsGetJobs(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	//TODO
	w.WriteHeader(http.StatusNotImplemented)
}

func (d *Server) scientistsApplyJob(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	//TODO
	w.WriteHeader(http.StatusNotImplemented)
}

func (d *Server) scientistsGetWorkspace(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	//TODO
	w.WriteHeader(http.StatusNotImplemented)
}

func (d *Server) companiesUploadCode(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	//TODO
	w.WriteHeader(http.StatusNotImplemented)
}

func (d *Server) execR(
	w http.ResponseWriter,
	req *http.Request,
	_ httprouter.Params,
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
	defer func() {
		if err := os.Chdir(cwd); err != nil {
			d.log.Printf("chdir: unexpected error %q", err)
		}
	}()

	cmd := exec.Command("R", "-f", "./code.r")
	d.log.Printf("executing R code")
	res, err := cmd.CombinedOutput()

	if err != nil {
		d.failrequest(w, "exec R: unexpected error %q", err)
		return
	}
	d.log.Printf("executed R code with success")

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	if err != nil {
		d.log.Printf("unexpected error %q sending response", err)
		return
	}
}

func (d *Server) receiveUpload(
	req *http.Request,
	filename string,
	jobID string,
) bool {
	uploadedfile, _, err := req.FormFile(filename)
	if err != nil {
		d.log.Printf("%q parsing form", err)
		return false
	}

	defer func() {
		if err := uploadedfile.Close(); err != nil {
			d.log.Printf("close: unexpected error %q", err)
		}
	}()

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

func (d *Server) failrequest(
	w http.ResponseWriter,
	fmt string,
	args ...interface{},
) {
	d.log.Printf(fmt, args...)
	w.WriteHeader(http.StatusInternalServerError)
}

// Handler ...
func (d *Server) Handler() http.Handler {
	return d.router
}

// ListenAndServe ...
func (d *Server) ListenAndServe(addr string) error {
	d.log.Printf("WebServer running at %q", addr)
	return http.ListenAndServe(addr, d.router)
}
