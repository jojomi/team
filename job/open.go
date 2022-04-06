package job

import (
	"fmt"
	"github.com/jojomi/go-script/v2"
	"net/url"
	"os"
	"runtime"

	"github.com/jojomi/go-script/v2/interview"
	"github.com/juju/errors"
	"github.com/pkg/browser"
	"gopkg.in/yaml.v2"
)

type OpenJob struct {
	Meta Metadata `yaml:"metadata"`

	Wait bool `yaml:"wait,omitempty"`

	Targets []string `yaml:"targets"`
}

func (x *OpenJob) Metadata() Metadata {
	return x.Meta
}

func (x *OpenJob) WithFilename(filename string) Job {
	x.Meta.Filename = filename
	return x
}

func (x *OpenJob) Pre() error {
	return nil
}

func (x *OpenJob) Post() error {
	return nil
}

func LoadOpenJobFromFile(filename string) (*OpenJob, error) {
	var x *OpenJob
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, errors.Annotatef(err, "could not load OpenJob from file %s", filename)
	}
	err = yaml.UnmarshalStrict(content, &x)
	if err != nil {
		return nil, errors.Annotatef(err, "invalid OpenJob configuration in file %s", filename)
	}

	x.Meta.Filename = filename
	return x, nil
}

func (x OpenJob) IsPossible() (bool, error) {
	return true, nil
}

func (x *OpenJob) Execute(options ExecutionOptions) error {
	for _, t := range x.Targets {
		if x.isValidUrl(t) {
			if options.DryRun {
				PrintDryRunTag()
				fmt.Printf("opening URL %s...\n", t)
				continue
			}
			err := openURLDefaultBrowser(t)
			if err != nil {
				return err
			}
			continue
		}

		if isBinary, runPath := x.isValidBinary(t); isBinary {
			if options.DryRun {
				PrintDryRunTag()
				fmt.Printf("opening binary %s...\n", t)
				continue
			}
			sc := script.NewContext()
			l := script.NewLocalCommand()
			l.Add(runPath)
			_, err := sc.ExecuteDetachedFullySilent(l)
			if err != nil {
				return err
			}
			continue
		}

		if options.DryRun {
			PrintDryRunTag()
			fmt.Printf("opening file %s...\n", t)
			continue
		}
		err := browser.OpenFile(t)
		if err != nil {
			return err
		}
	}

	if x.Wait {
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

func (x OpenJob) isValidBinary(toTest string) (isBinary bool, runPath string) {
	sc := script.NewContext()
	if sc.CommandExists(toTest) {
		return true, sc.CommandPath(toTest)
	}

	// TODO: check if file is executable itself (with PATH or absolut path)
	// see https://gitlab.com/stellarpower-grouped-projects/tidbits/go/-/blob/main/CheckFileExecutable.go#L10

	return false, ""
}

func (x OpenJob) isValidUrl(toTest string) bool {
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

// open opens the specified URL in the default browser of the user.
func openURLDefaultBrowser(url string) error {
	l := script.NewLocalCommand()

	switch runtime.GOOS {
	case "windows":
		l.AddAll("cmd", "/c", "start")
	case "darwin":
		l.Add("open")
	default: // "linux", "freebsd", "openbsd", "netbsd"
		l.Add("xdg-open")
	}
	l.Add(url)

	sc := script.NewContext()
	pr, err := sc.ExecuteFullySilent(l)

	if err != nil {
		return err
	}
	if !pr.Successful() {
		return fmt.Errorf("could not open url %s", url)
	}

	return nil
}
