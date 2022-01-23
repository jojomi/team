package main

import (
	"errors"
	"github.com/spf13/cobra"
)

// EnvAssert encapsulates the environment for the CLI assert handler.
type EnvAssert struct {
	Command string
	Args    []string

	Verbose bool
}

// ParseFrom reads the state from a given cobra command and its args.
func (e *EnvAssert) ParseFrom(command *cobra.Command, args []string) error {
	var err error

	if len(args) < 1 {
		return errors.New("no command given")
	}
	e.Command = args[0]

	if len(args) > 1 {
		e.Args = args[1:]
	}

	e.Verbose, err = command.Flags().GetBool("verbose")
	if err != nil {
		return err
	}
	return nil
}
