package main

import (
	"fmt"
	"github.com/jojomi/team/assert"
	"github.com/jojomi/team/exit"
	"math/rand"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func getAssertCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "assert",
		Run: handleAssertCmd,
	}

	return cmd
}

func handleAssertCmd(cmd *cobra.Command, args []string) {
	var exitCode exit.Code
	env := EnvAssert{}
	err := env.ParseFrom(cmd, args)
	if err != nil {
		log.Fatal().Err(err).Msg("could not parse command")
	}
	err, exitCode = handleAssert(env)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	os.Exit(int(exitCode))
}

func handleAssert(env EnvAssert) (error, exit.Code) {
	setLoggerVerbosity(env.Verbose)

	handlerMap := map[string]func([]string) (error, exit.Code){
		"command-exists":               assert.CommandExists,
		"file-exists":                  assert.FileExists,
		"dir-exists":                   assert.DirExists,
		"mounted":                      assert.Mounted,
		"non-empty-dir":                assert.NonEmptyDir,
		"ssh-reachable":                assert.SSHReachable,
		"ssh-reachable-noninteractive": assert.SSHReachableNonInteractive,
		"weekday":                      assert.Weekday,
	}

	cmd := env.Command
	f, ok := handlerMap[cmd]
	if !ok {
		return fmt.Errorf("assert handler not found: %s", cmd), exit.CodeErrorFinal
	}

	return f(env.Args)
}
