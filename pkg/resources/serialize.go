package resources

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
)

// ParseJson parses JSON from the file in embedded resources
func ParseJson(path string, out interface{}) error {
	data, err := ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, out)
}

// ParseYaml parses yaml from the file in embedded resources
func ParseYaml(path string, out interface{}) error {
	data, err := ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, out)
}
