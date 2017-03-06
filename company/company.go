package company

import (
	"encoding/json"
	"io/ioutil"
)

// Job ...
type Job struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Status      string  `json:"status"`
	LastUpdate  string  `json:"lastUpdate"`
}

// Company ...
type Company struct {
	jobs []Job
}

// NewCompany ...
func NewCompany() *Company {
	jobs := loadJobs()
	return &Company{jobs: jobs}
}

func loadJobs() []Job {
	raw, err := ioutil.ReadFile("./db/company/jobs.json")
	if err != nil {
		return nil
	}

	var jobs []Job
	if err := json.Unmarshal(raw, &jobs); err != nil {
		return nil
	}
	return jobs
}

// GetJobsByStatus ...
func (c *Company) GetJobsByStatus(status string) []Job {
	jobs := make([]Job, 0, len(c.jobs))

	for _, job := range c.jobs {
		if job.Status == status {
			jobs = append(jobs, job)
		}
	}

	return jobs
}

// AddNewJob ...
func (c *Company) AddNewJob(job Job) {
	job.ID = len(c.jobs)
	job.Status = "pending"

	c.jobs = append(c.jobs, job)
}

// GetJob ...
func (c *Company) GetJob(id int) *Job {
	if id >= len(c.jobs) {
		return nil
	}

	return &c.jobs[id]
}
