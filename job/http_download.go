package job

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/alexeyco/goozzle"
	"github.com/juju/errors"
	"gopkg.in/yaml.v2"
)

type HTTPDownloadJob struct {
	Meta Metadata `yaml:"metadata"`

	Files []struct {
		RemoteURL     string `yaml:"remote-url"`
		LocalDir      string `yaml:"local-dir"`
		LocalFilename string `yaml:"local-filename"`
	} `yaml:"files"`
}

func (x *HTTPDownloadJob) Metadata() Metadata {
	return x.Meta
}

func (x *HTTPDownloadJob) WithFilename(filename string) *HTTPDownloadJob {
	x.Meta.Filename = filename
	return x
}

func (x *HTTPDownloadJob) Pre() error {
	return nil
}

func (x *HTTPDownloadJob) Post() error {
	return nil
}

func LoadHTTPDownloadJobFromFile(filename string) (*HTTPDownloadJob, error) {
	var x *HTTPDownloadJob
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, errors.Annotatef(err, "could not load HTTPDownloadJob from file %s", filename)
	}
	err = yaml.UnmarshalStrict(content, &x)
	if err != nil {
		return nil, errors.Annotatef(err, "invalid HTTPDownloadJob configuration in file %s", filename)
	}

	x.Meta.Filename = filename
	return x, nil
}

func (x HTTPDownloadJob) IsPossible() (bool, error) {
	// at least one URL reachable?
	var url string
	for _, file := range x.Files {
		url = file.RemoteURL
		if isURLReachable(url) {
			return true, nil
		}
	}

	return false, fmt.Errorf("URL %s not reachable", url)
}

func (x *HTTPDownloadJob) Execute(options ExecutionOptions) error {
	var (
		localTargetFile string
	)

	for _, file := range x.Files {
		localTargetFile = filepath.Join(file.LocalDir, file.LocalFilename)

		// dry run?
		if options.DryRun {
			PrintDryRunTag()
			fmt.Printf("Downloading %s to %s...\n", file.RemoteURL, localTargetFile)
			continue
		}

		// actually download
		u, _ := url.Parse(file.RemoteURL)

		res, err := goozzle.Get(u).Do()
		if err != nil {
			return errors.Annotatef(err, "Error downloading %s", file.RemoteURL)
		}

		os.WriteFile(localTargetFile, res.Body(), 0750)
	}
	fmt.Printf("downloaded %d URLs\n", len(x.Files))
	return nil
}
