package datahub

import (
	"encoding/json"
	"io/ioutil"
)

// Job ...
type Job struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
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
	raw, err := ioutil.ReadFile("./docs/companies_jobs.json")
	if err != nil {
		return nil
	}

	var jobs []Job
	if err := json.Unmarshal(raw, &jobs); err != nil {
		return nil
	}
	return jobs
}

// GetJobs ...
func (c *Company) GetJobs() []Job {
	return c.jobs
}
