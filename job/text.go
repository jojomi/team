package job

import (
	"fmt"
	"net/url"
	"os"

	"github.com/jojomi/go-script/v2/interview"
	"github.com/juju/errors"
	"gopkg.in/yaml.v2"
)

type TextJob struct {
	Meta Metadata `yaml:"metadata"`

	Wait bool `yaml:"wait,omitempty"`

	Text string `yaml:"text,omitempty"`
}

func (x *TextJob) Metadata() Metadata {
	return x.Meta
}

func (x *TextJob) WithFilename(filename string) Job {
	x.Meta.Filename = filename
	return x
}

func (x *TextJob) Pre() error {
	return execCommand(x.Meta.Pre.Command)
}

func (x *TextJob) Post() error {
	return execCommand(x.Meta.Post.Command)
}

func LoadTextJobFromFile(filename string) (*TextJob, error) {
	var x *TextJob
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, errors.Annotatef(err, "could not load TextJob from file %s", filename)
	}
	err = yaml.UnmarshalStrict(content, &x)
	if err != nil {
		return nil, errors.Annotatef(err, "invalid TextJob configuration in file %s", filename)
	}

	x.Meta.Filename = filename
	return x, nil
}

func (x TextJob) IsPossible() (bool, error) {
	return true, nil
}

func (x *TextJob) Execute(options ExecutionOptions) error {

	fmt.Println("")
	fmt.Println(x.Text)

	if x.Wait && options.Wait {
		done, err := interview.Confirm("Task done?", false)
		if err != nil {
			return err
		}
		if !done {
			return errors.New("task not done")
		}
	}

	return nil
}

func (x TextJob) isValidUrl(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}
