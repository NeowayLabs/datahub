package scientists

import (
	"encoding/json"
	"io/ioutil"
)

// Scientist ...
type Scientist struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	Rating int    `json:"rating"`
	Likes  int    `json:"likes"`
}

// Scientists ...
type Scientists struct {
	scientists []*Scientist
}

// NewScientists ...
func NewScientists() *Scientists {
	scientists := loadScientists()
	return &Scientists{scientists: scientists}
}

func loadScientists() []*Scientist {
	raw, err := ioutil.ReadFile("./db/scientists/scientists.json")
	if err != nil {
		return nil
	}

	var scientists []*Scientist
	if err := json.Unmarshal(raw, &scientists); err != nil {
		return nil
	}
	return scientists
}

// GetScientists ...
func (s *Scientists) GetScientists() []*Scientist {
	return s.scientists
}
