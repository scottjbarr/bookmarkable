package bookmarkable

import (
	"encoding/json"
	"io/ioutil"
)

// Config top level configuration struct
//
// GistName string `json:"gist_name"`
type config struct {
	ID    *string `json:"id"`
	Token string  `json:"token"`
}

// ParseConfig parses a config file
func parseConfig(filename string) (*config, error) {
	bytes, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	config := config{}

	if err := json.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
