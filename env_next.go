package main

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// EnvNext encapsulates the environment for the CLI next handler.
type EnvNext struct {
	JobDir         string
	UnattendedOnly bool
	PossibleOnly   bool
	Fixable        bool
	ShowUnfixable  bool

	JobFile string

	Verbose bool
}

// ParseFrom reads the state from a given cobra command and its args.
func (e *EnvNext) ParseFrom(command *cobra.Command, _ []string) error {
	var err error

	jobDir, err := command.Flags().GetString("job-dir")
	if err != nil {
		return err
	}
	e.JobDir, err = homedir.Expand(jobDir)
	if err != nil {
		return err
	}

	e.UnattendedOnly, err = command.Flags().GetBool("unattended-only")
	if err != nil {
		return err
	}
	e.PossibleOnly, err = command.Flags().GetBool("possible-only")
	if err != nil {
		return err
	}
	e.Fixable, err = command.Flags().GetBool("fixable")
	if err != nil {
		return err
	}

	e.ShowUnfixable, err = command.Flags().GetBool("show-unfixable")
	if err != nil {
		return err
	}

	e.JobFile, err = command.Flags().GetString("job-file")
	if err != nil {
		return err
	}

	e.Verbose, err = command.Flags().GetBool("verbose")
	if err != nil {
		return err
	}

	return nil
}
