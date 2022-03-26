package main

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// EnvHealth encapsulates the environment for the CLI health handler.
type EnvHealth struct {
	JobDir  string
	JobFile string
	Command string
	Args    []string

	RequiredOnly bool
	PossibleOnly bool

	Verbose bool
}

// ParseFrom reads the state from a given cobra command and its args.
func (e *EnvHealth) ParseFrom(command *cobra.Command, args []string) error {
	var err error

	jobDir, err := command.Flags().GetString("job-dir")
	if err != nil {
		return err
	}
	e.JobDir, err = homedir.Expand(jobDir)
	if err != nil {
		return err
	}

	e.JobFile, err = command.Flags().GetString("job-file")
	if err != nil {
		return err
	}

	e.RequiredOnly, err = command.Flags().GetBool("required-only")
	if err != nil {
		return err
	}
	e.PossibleOnly, err = command.Flags().GetBool("possible-only")
	if err != nil {
		return err
	}

	e.Verbose, err = command.Flags().GetBool("verbose")
	if err != nil {
		return err
	}
	return nil
}
