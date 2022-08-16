package main

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// EnvRoot encapsulates the environment for the CLI root handler.
type EnvRoot struct {
	JobDir         string
	Manual         bool
	SkipTimeCheck  bool
	DefaultYes     bool
	UnattendedOnly bool
	AttendedOnly   bool
	PossibleOnly   bool
	Fixable        bool
	Shutdown       bool
	SelectFirst    bool
	ShowUnfixable  bool

	JobFile string

	DryRun  bool
	Verbose bool
}

// ParseFrom reads the state from a given cobra command and its args.
func (e *EnvRoot) ParseFrom(command *cobra.Command, _ []string) error {
	var err error

	jobDir, err := command.Flags().GetString("job-dir")
	if err != nil {
		return err
	}
	e.JobDir, err = homedir.Expand(jobDir)
	if err != nil {
		return err
	}

	e.Manual, err = command.Flags().GetBool("manual")
	if err != nil {
		return err
	}
	e.SelectFirst, err = command.Flags().GetBool("select-first")
	if err != nil {
		return err
	}

	e.DefaultYes, err = command.Flags().GetBool("default-yes")
	if err != nil {
		return err
	}
	e.SkipTimeCheck, err = command.Flags().GetBool("skip-time-check")
	if err != nil {
		return err
	}

	e.UnattendedOnly, err = command.Flags().GetBool("unattended-only")
	if err != nil {
		return err
	}
	e.AttendedOnly, err = command.Flags().GetBool("attended-only")
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

	e.Shutdown, err = command.Flags().GetBool("shutdown")
	if err != nil {
		return err
	}

	e.JobFile, err = command.Flags().GetString("job-file")
	if err != nil {
		return err
	}

	e.DryRun, err = command.Flags().GetBool("dry-run")
	if err != nil {
		return err
	}
	e.Verbose, err = command.Flags().GetBool("verbose")
	if err != nil {
		return err
	}
	return nil
}

func (e *EnvRoot) String() string {
	out, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("EnvRoot %s", string(out))
}
