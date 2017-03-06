package company

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Scientist ...
type Scientist struct {
	ID              int     `json:"id"`
	Name            string  `json:"name"`
	Rating          int     `json:"rating"`
	Counterproposal float64 `json:"counterproposal,omitempty"`
}

// Job ...
type Job struct {
	ID          int          `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Price       float64      `json:"price"`
	Status      string       `json:"status"`
	LastUpdate  string       `json:"lastUpdate"`
	Candidates  []*Scientist `json:"candidates,omitempty"`
	Scientists  []*Scientist `json:"scientists,omitempty"`
}

// Company ...
type Company struct {
	jobs []*Job
}

// NewCompany ...
func NewCompany() *Company {
	jobs := loadJobs()
	return &Company{jobs: jobs}
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

// GetJobsByStatus ...
func (c *Company) GetJobsByStatus(status string) []*Job {
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
	job.ID = len(c.jobs)
	job.Status = "pending"

	c.jobs = append(c.jobs, job)
}

// GetJob ...
func (c *Company) GetJob(id int) *Job {
	if id >= len(c.jobs) {
		return nil
	}

	return c.jobs[id]
}

// StartJob ...
func (c *Company) StartJob(id int, scientists []*Scientist) error {
	if id >= len(c.jobs) {
		return fmt.Errorf("Job %d not found", id)
	}

	job := c.jobs[id]

	job.Scientists = make([]*Scientist, 0, len(scientists))

	for _, scientist := range scientists {
		for _, candidate := range job.Candidates {
			if candidate.ID == scientist.ID {
				job.Scientists = append(job.Scientists, candidate)
				break
			}
		}
	}

	job.Candidates = nil
	job.Status = "doing"

	return nil
}

// ApplyScientist ...
func (c *Company) ApplyScientist(id int, candidate *Scientist) error {
	if id >= len(c.jobs) {
		return fmt.Errorf("Job %d not found", id)
	}

	job := c.jobs[id]
	job.Candidates = append(job.Candidates, candidate)
	return nil
}
