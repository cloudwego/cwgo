package docker

import (
	"gopkg.in/yaml.v3"
	"os"
)

type DockerfilesTpl struct {
	Body string `json:"body"`
	Path string `json:"path"`
}

// FromYAMLFile unmarshals a DockerfilesTpl with YAML format from the given file.
func (p *DockerfilesTpl) FromYAMLFile(filename string) error {
	if p == nil {
		return nil
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, p)
}
