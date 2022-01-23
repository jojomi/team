package job

import (
	"fmt"
	"os"

	"github.com/juju/errors"
	"gopkg.in/yaml.v2"
)

type Job interface {
	Metadata() Metadata

	// Execution
	IsPossible() (bool, error)
	Pre() error
	Execute(options ExecutionOptions) error
	Post() error
}

type Factory struct{}

func (x Factory) GetFromFile(filename string) (Job, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = yaml.Unmarshal(content, &data)
	if err != nil {
		return nil, err
	}
	metadata, ok := data["metadata"]
	if !ok {
		return nil, errors.New("metadata.type not set")
	}
	jobType, ok := metadata.(map[interface{}]interface{})["type"]
	if !ok {
		return nil, errors.New("metadata.type not set")
	}

	switch jobType {
	case "feed-download":
		return LoadFeedDownloadJobFromFile(filename)
	case "http-download":
		return LoadHTTPDownloadJobFromFile(filename)
	case "shell":
		return LoadShellJobFromFile(filename)
	case "open":
		return LoadOpenJobFromFile(filename)
	case "text":
		return LoadTextJobFromFile(filename)
	}
	return nil, fmt.Errorf("could not find job with type %s", jobType)
}
