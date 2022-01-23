package job

import (
	"fmt"
	"github.com/jojomi/go-script/v2/interview"
	"os"

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
	return nil
}

func (x *ShellJob) Post() error {
	return nil
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
	c := x.getExecuteCommand()
	if !sc.CommandExists(c.Binary()) && !sc.FileExists(c.Binary()) {
		return false, fmt.Errorf("missing command: %s", c.Binary())
	}

	return true, nil
}

func (x *ShellJob) getExecuteCommand() script.Command {
	command, err := homedir.Expand(x.Exec)
	if err != nil {
		panic(err)
	}
	return script.LocalCommandFrom(command)
}

func (x *ShellJob) getExecuteDryRunCommand() script.Command {
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
	cmd := x.getExecuteCommand()

	executor := getExecutorByOutputType(x.Meta.Output)

	debugCommand(cmd, options)

	if options.DryRun {
		if cmd = x.getExecuteDryRunCommand(); cmd == nil {
			return nil
		}
	}

	pr, err := executor(cmd)
	if err != nil {
		return err
	}
	if !pr.Successful() {
		return fmt.Errorf("unsuccessful command: %s", cmd.String())
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
