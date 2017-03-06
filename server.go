package datahub

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/NeowayLabs/datahub/company"
	"github.com/NeowayLabs/datahub/scientists"
	"github.com/julienschmidt/httprouter"
)

// Server ...
type Server struct {
	router     *httprouter.Router
	datadir    string
	log        *log.Logger
	company    *company.Company
	scientists *scientists.Scientists
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
	scientists := scientists.NewScientists()

	d := &Server{
		router:     router,
		datadir:    datadir,
		log:        log,
		company:    company,
		scientists: scientists,
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
	New     []*company.Job `json:"new,omitempty"`
	Pending []*company.Job `json:"pending"`
	Doing   []*company.Job `json:"doing"`
	Done    []*company.Job `json:"done"`
}

func (d *Server) companiesGetJobs(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	defer func() {
		if err := req.Body.Close(); err != nil {
			d.log.Printf("body close: error %q", err)
		}
	}()

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

func (d *Server) companiesCreateJob(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	defer func() {
		if err := req.Body.Close(); err != nil {
			d.log.Printf("body close: error %q", err)
		}
	}()

	decoder := json.NewDecoder(req.Body)

	var job company.Job
	if err := decoder.Decode(&job); err != nil {
		d.log.Printf("unmarshal: error %q", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	d.company.AddNewJob(&job)

	w.WriteHeader(http.StatusOK)
}

func getID(params httprouter.Params, name string) (int, error) {
	s := params.ByName(name)
	if s == "" {
		return 0, fmt.Errorf("param: id is empty")
	}

	id, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("param: id is not a number")
	}

	return id, nil
}

func (d *Server) companiesGetJob(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	defer func() {
		if err := req.Body.Close(); err != nil {
			d.log.Printf("body close: error %q", err)
		}
	}()

	id, err := getID(params, "id")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	job := d.company.GetJob(id)
	if job == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	bytes, err := json.Marshal(job)
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

func (d *Server) companiesStartJob(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	defer func() {
		if err := req.Body.Close(); err != nil {
			d.log.Printf("body close: error %q", err)
		}
	}()

	type Scientists struct {
		Scientists []*company.Scientist `json:"scientists"`
	}

	job, err := getID(params, "id")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(req.Body)

	var scientists Scientists
	if err := decoder.Decode(&scientists); err != nil {
		d.log.Printf("unmarshal: error %q", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := d.company.StartJob(job, scientists.Scientists); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (d *Server) scientistsList(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	defer func() {
		if err := req.Body.Close(); err != nil {
			d.log.Printf("body close: error %q", err)
		}
	}()

	scientists := d.scientists.GetScientists()

	bytes, err := json.Marshal(scientists)
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

func (d *Server) scientistsGetJobs(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	defer func() {
		if err := req.Body.Close(); err != nil {
			d.log.Printf("body close: error %q", err)
		}
	}()

	id, err := getID(params, "id")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pendingJobs := d.company.GetJobsByStatus("pending")
	doingJobs := d.company.GetJobsByStatus("doing")
	doneJobs := d.company.GetJobsByStatus("done")

	new := make([]*company.Job, 0, len(pendingJobs))
	pending := make([]*company.Job, 0, len(pendingJobs))
	doing := make([]*company.Job, 0, len(doingJobs))
	done := make([]*company.Job, 0, len(doneJobs))

jobs:
	for _, job := range pendingJobs {
		candidates := job.Candidates
		for _, candidate := range candidates {
			if candidate.ID == id {
				pending = append(pending, job)
				continue jobs
			}
		}
		new = append(new, job)
	}

	for _, job := range doingJobs {
		scientists := job.Scientists
		for _, scientist := range scientists {
			if scientist.ID == id {
				doing = append(doing, job)
				break
			}
		}
	}

	for _, job := range doneJobs {
		scientists := job.Scientists
		for _, scientist := range scientists {
			if scientist.ID == id {
				done = append(done, job)
				break
			}
		}
	}

	jobs := &Jobs{
		New:     new,
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

// Apply ...
type Apply struct {
	Counterproposal float64 `json:"counterproposal"`
}

func (d *Server) scientistsApplyJob(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	defer func() {
		if err := req.Body.Close(); err != nil {
			d.log.Printf("body close: error %q", err)
		}
	}()

	id, err := getID(params, "id")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	scientist := d.scientists.GetScientist(id)

	job, err := getID(params, "job")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(req.Body)

	var apply Apply
	if err := decoder.Decode(&apply); err != nil {
		d.log.Printf("unmarshal: error %q", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	candidate := &company.Scientist{
		ID:              scientist.ID,
		Name:            scientist.Name,
		Rating:          scientist.Rating,
		Counterproposal: apply.Counterproposal,
	}

	if err := d.company.ApplyScientist(job, candidate); err != nil {
		d.log.Printf("apply: error %q", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (d *Server) scientistsGetWorkspace(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Workspace ...
	type Workspace struct {
		ID               int               `json:"id"`
		Title            string            `json:"title"`
		Description      string            `json:"description"`
		Proposed         float64           `json:"proposed"`
		AccuracyRequired float64           `json:"accuracyRequired"`
		Deadline         string            `json:"deadline"`
		Status           string            `json:"status"`
		LastUpdate       string            `json:"lastUpdate"`
		Workspace        company.Workspace `json:"workspace"`
	}

	defer func() {
		if err := req.Body.Close(); err != nil {
			d.log.Printf("body close: error %q", err)
		}
	}()

	scientistID, err := getID(params, "id")
	if err != nil {
		d.log.Printf("params: %q", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	jobID, err := getID(params, "job")
	if err != nil {
		d.log.Printf("params: %q", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	job := d.company.GetJob(jobID)
	if job == nil {
		d.log.Printf("job %d not found", jobID)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var scientist *company.Scientist

	for _, sc := range job.Scientists {
		if sc.ID == scientistID {
			scientist = sc
			break
		}
	}

	if scientist == nil {
		d.log.Printf("scientist %d not found", scientistID)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	workspace := &Workspace{
		ID:               job.ID,
		Title:            job.Title,
		Description:      job.Description,
		Proposed:         job.Proposed,
		AccuracyRequired: job.AccuracyRequired,
		Deadline:         job.Deadline,
		Status:           job.Status,
		LastUpdate:       job.LastUpdate,
		Workspace:        scientist.Workspace,
	}

	bytes, err := json.Marshal(workspace)
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

func (d *Server) companiesUploadCode(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	//TODO executar o codigo em R e atualizar os campos do workspace do scientist: accuracy, result, code

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

func (d *Server) scorecheck(predictionfilepath string, resultfilepath string) (float32, error) {

	predictionfile, err := os.Open(predictionfilepath)
	if err != nil {
		d.log.Printf("error: %q reading file", err)
		return 0, err
	}
	defer predictionfile.Close()

	resultfile, err := os.Open(resultfilepath)
	if err != nil {
		d.log.Printf("error: %q reading file", err)
		return 0, err
	}
	defer resultfile.Close()

	scanpredictionfile := bufio.NewScanner(predictionfile)
	scanresultfile := bufio.NewScanner(resultfile)

	var ok, totallines float32

	for scanresultfile.Scan() {
		d.log.Printf("reading result line %q", scanresultfile.Text())
		totallines = totallines + 1
		if scanpredictionfile.Scan() {
			d.log.Printf("reading prediction line %q", scanpredictionfile.Text())
			if scanresultfile.Text() == scanpredictionfile.Text() {
				ok++
				d.log.Printf("line equal! %q", scanpredictionfile.Text())
			} else {
				d.log.Printf("line NOT equal! %q", scanpredictionfile.Text())
			}
		}
	}
	score := float32(int((ok*100/totallines)*100)) / 100
	d.log.Printf("detected score: %f", score)

	return score, nil
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
