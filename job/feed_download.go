package job

import (
	"fmt"
	"os"

	"github.com/jojomi/go-script/v2"
	"github.com/juju/errors"
	"gopkg.in/yaml.v2"
)

type FeedDownloadJob struct {
	Meta Metadata `yaml:"metadata"`

	RemoteURL string `yaml:"remote-url"`
	LocalDir  string `yaml:"local-dir"`
}

func (x *FeedDownloadJob) Metadata() Metadata {
	return x.Meta
}

func (x *FeedDownloadJob) WithFilename(filename string) *FeedDownloadJob {
	x.Meta.Filename = filename
	return x
}

func (x *FeedDownloadJob) Pre() error {
	return nil
}

func (x *FeedDownloadJob) Post() error {
	return nil
}

func LoadFeedDownloadJobFromFile(filename string) (*FeedDownloadJob, error) {
	var x *FeedDownloadJob
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, errors.Annotatef(err, "could not load FeedDownloadJob from file %s", filename)
	}
	err = yaml.UnmarshalStrict(content, &x)
	if err != nil {
		return nil, errors.Annotatef(err, "invalid FeedDownloadJob configuration in file %s", filename)
	}
	x.Meta.Filename = filename
	return x, nil
}

func (x FeedDownloadJob) IsPossible() (bool, error) {
	sc := script.NewContext()

	// feeddownload tool installed?
	if !sc.CommandExists("feeddownload") {
		return false, fmt.Errorf("feeddownload not found in PATH, get it from https://github.com/jojomi/feeddownload")
	}

	// URL reachable?
	url := x.RemoteURL
	if !isURLReachable(url) {
		return false, fmt.Errorf("URL %s not reachable", url)
	}

	return true, nil
}

func (x *FeedDownloadJob) Execute(options ExecutionOptions) error {
	cm := script.NewLocalCommand()
	cm.AddAll("feeddownload", x.RemoteURL, x.LocalDir)

	// dry run?
	if options.DryRun {
		cm.Add("--dry-run")
	}

	debugCommand(cm, options)

	executor := getExecutorByOutputType(x.Metadata().Output, nil)

	pr, err := executor(cm)
	if err != nil {
		return errors.Annotatef(err, "could not execute successfully %s", cm.String())
	}
	if !pr.Successful() {
		return errors.New("error while executing feeddownload")
	}
	return nil
}
