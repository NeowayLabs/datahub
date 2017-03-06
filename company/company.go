package company

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
)

// Workspace ...
type Workspace struct {
	Accuracy    int     `json:"accuracy"`
	Result      string  `json:"result"`
	Code        int     `json:"code"`
	Description float64 `json:"Description"`
}

// Scientist ...
type Scientist struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	Rating          int       `json:"rating"`
	Counterproposal float64   `json:"counterproposal,omitempty"`
	Workspace       Workspace `json:"workspace"`
}

// Job ...
type Job struct {
	ID               int          `json:"id"`
	Title            string       `json:"title"`
	Description      string       `json:"description"`
	Proposed         float64      `json:"proposed"`
	AccuracyRequired float64      `json:"accuracyRequired"`
	Deadline         string       `json:"deadline"`
	Status           string       `json:"status"`
	LastUpdate       string       `json:"lastUpdate"`
	Candidates       []*Scientist `json:"candidates,omitempty"`
	Scientists       []*Scientist `json:"scientists,omitempty"`
}

// Company ...
type Company struct {
	mutex *sync.RWMutex
	jobs  []*Job
}

// NewCompany ...
func NewCompany() *Company {
	jobs := loadJobs()
	mutex := &sync.RWMutex{}

	return &Company{
		jobs:  jobs,
		mutex: mutex,
	}
}

func loadJobs() []*Job {
	raw, err := ioutil.ReadFile("./db/company/jobs.json")
	if err != nil {
		return nil
	}

	var jobs []*Job
	if err := json.Unmarshal(raw, &jobs); err != nil {
		return nil
	}
	return jobs
}

func saveJobs(jobs []*Job) {
	buffer, err := json.Marshal(jobs)
	if err != nil {
		return
	}

	err = ioutil.WriteFile("./db/company/jobs.json", buffer, 0644)
	if err != nil {
		return
	}
}

// GetJobsByStatus ...
func (c *Company) GetJobsByStatus(status string) []*Job {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	jobs := make([]*Job, 0, len(c.jobs))

	for _, job := range c.jobs {
		if job.Status == status {
			jobs = append(jobs, job)
		}
	}

	return jobs
}

// AddNewJob ...
func (c *Company) AddNewJob(job *Job) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	job.ID = len(c.jobs)
	job.Status = "pending"

	c.jobs = append(c.jobs, job)

	saveJobs(c.jobs)
}

// GetJob ...
func (c *Company) GetJob(id int) *Job {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if id >= len(c.jobs) {
		return nil
	}

	return c.jobs[id]
}

// StartJob ...
func (c *Company) StartJob(id int, scientists []*Scientist) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if id >= len(c.jobs) {
		return fmt.Errorf("Job %d not found", id)
	}

	job := c.jobs[id]

	job.Scientists = make([]*Scientist, 0, len(scientists))

	for _, scientist := range scientists {
		for _, candidate := range job.Candidates {
			if candidate.ID == scientist.ID {
				//candidate.Workspace = Workspace{}
				job.Scientists = append(job.Scientists, candidate)
				break
			}
		}
	}

	job.Candidates = nil
	job.Status = "doing"

	saveJobs(c.jobs)
	return nil
}

// ApplyScientist ...
func (c *Company) ApplyScientist(id int, candidate *Scientist) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if id >= len(c.jobs) {
		return fmt.Errorf("Job %d not found", id)
	}

	job := c.jobs[id]
	job.Candidates = append(job.Candidates, candidate)

	saveJobs(c.jobs)
	return nil
}
