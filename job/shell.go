package job

import (
	"fmt"
	"github.com/jojomi/gorun"
	"os"
	"path"

	"github.com/jojomi/go-script/v2/interview"

	"github.com/jojomi/go-script/v2"
	"github.com/juju/errors"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

type ShellJob struct {
	Meta Metadata `yaml:"metadata"`

	Wait bool `yaml:"wait"`

	Exec       string `yaml:"execute"`
	ExecDryRun string `yaml:"execute-dry-run,omitempty"`
}

func (x *ShellJob) Metadata() Metadata {
	return x.Meta
}

func (x *ShellJob) WithFilename(filename string) Job {
	x.Meta.Filename = filename
	return x
}

func (x *ShellJob) Pre() error {
	return execCommand(x.Meta.Pre.Command)
}

func (x *ShellJob) Post() error {
	return execCommand(x.Meta.Post.Command)
}

func LoadShellJobFromFile(filename string) (*ShellJob, error) {
	var x *ShellJob
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, errors.Annotatef(err, "could not load ShellJob from file %s", filename)
	}
	err = yaml.UnmarshalStrict(content, &x)
	if err != nil {
		return nil, errors.Annotatef(err, "invalid ShellJob configuration in file %s", filename)
	}

	x.Meta.Filename = filename
	return x, nil
}

func (x ShellJob) IsPossible() (bool, error) {
	// binary available?
	sc := script.NewContext()
	sc.SetWorkingDir(path.Dir(x.Meta.Filename))
	c := x.getExecuteCommand()
	if !sc.CommandExists(c.Binary()) && !sc.FileExists(c.Binary()) {
		return false, fmt.Errorf("missing command: %s, working dir is %s", c.Binary(), sc.WorkingDir())
	}

	return true, nil
}

func (x *ShellJob) getExecuteCommand() gorun.Command {
	command, err := homedir.Expand(x.Exec)
	if err != nil {
		panic(err)
	}
	return gorun.LocalCommandFrom(command)
}

func (x *ShellJob) getExecuteDryRunCommand() gorun.Command {
	if x.ExecDryRun == "" {
		return nil
	}
	command, err := homedir.Expand(x.ExecDryRun)
	if err != nil {
		panic(err)
	}
	return script.LocalCommandFrom(command)
}

func (x *ShellJob) Execute(options ExecutionOptions) error {
	command := x.getExecuteCommand()
	debugGorunCommand(command, options)

	if options.DryRun {
		command = x.getExecuteDryRunCommand()
	}

	runner := gorun.NewWithCommand(command)
	runner.AddEnv("TEAM_DRYRUN", envBool(options.DryRun))
	runner.AddEnv("TEAM_VERBOSE", envBool(options.Verbose))
	runner.InWorkingDir(path.Dir(x.Meta.Filename))
	// runner.WithoutStdout()

	result, err := runner.Exec()
	if err != nil {
		return err
	}
	if !result.Successful() {
		return fmt.Errorf("unsuccessful command: %s", command.String())
	}

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

func envBool(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
