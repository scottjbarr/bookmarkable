package main

import (
	"encoding/json"
	"io/ioutil"
)

// Config top level configuration struct
//
// GistName string `json:"gist_name"`
type Config struct {
	Token string `json:"token"`
}

// ParseConfig parses a config file
func parseConfig(filename string) (*Config, error) {
	bytes, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	config := Config{}

	if err := json.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
