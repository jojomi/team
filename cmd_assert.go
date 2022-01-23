package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jojomi/go-script/v2"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func getAssertCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "assert",
		Run: handleAssertCmd,
	}

	return cmd
}

func handleAssertCmd(cmd *cobra.Command, args []string) {
	env := EnvAssert{}
	err := env.ParseFrom(cmd, args)
	if err != nil {
		log.Fatal().Err(err).Msg("could not parse command")
	}
	err = handleAssert(env)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func handleAssert(env EnvAssert) error {
	setLoggerVerbosity(env.Verbose)

	handlerMap := map[string]func([]string) error{
		"command-exists": assertCommandExists,
		"file-exists":    assertFileExists,
		"dir-exists":     assertDirExists,
		"mounted":        assertMounted,
		"non-empty-dir":  assertNonEmptyDir,
		"ssh-reachable":  assertSSHReachable,
	}

	cmd := env.Command
	f, ok := handlerMap[cmd]
	if !ok {
		return fmt.Errorf("assert handler not found: %s", cmd)
	}

	return f(env.Args)
}

func assertCommandExists(args []string) error {
	if len(args) < 1 {
		return errors.New("no command name given")
	}
	sc := script.NewContext()
	for _, cmd := range args {
		if sc.CommandExists(cmd) {
			continue
		}
		return fmt.Errorf("%s not installed or not in PATH", cmd)
	}
	return nil
}

func assertFileExists(args []string) error {
	if len(args) < 1 {
		return errors.New("no filename name given")
	}
	sc := script.NewContext()
	for _, filename := range args {
		if sc.FileExists(filename) {
			continue
		}
		return fmt.Errorf("file %s does not exist", filename)
	}
	return nil
}

func assertMounted(args []string) error {
	if len(args) < 1 {
		return errors.New("no mountpoint dir given")
	}
	sc := script.NewContext()
	for _, filename := range args {
		if !sc.DirExists(filename) {
			return fmt.Errorf("mounting dir %s does not exist", filename)
		}

	}
	return nil
}

func assertDirExists(args []string) error {
	sc := script.NewContext()

	for _, dir := range args {
		if !sc.DirExists(dir) {
			return fmt.Errorf("directory %s does not exist", dir)
		}
	}
	return nil
}

func assertNonEmptyDir(args []string) error {
	for _, dir := range args {
		if err := assertDirExists([]string{dir}); err != nil {
			return err
		}

		f, err := os.Open(dir)
		if err != nil {
			return err
		}

		_, err = f.Readdirnames(1)
		if err == io.EOF {
			return fmt.Errorf("directory %s if not empty", dir)
		}

		err = f.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func assertSSHReachable(args []string) error {
	sc := script.NewContext()
	var (
		l      *script.LocalCommand
		random = "abcabc"
	)
	for _, host := range args {
		l = script.NewLocalCommand()
		l.AddAll("ssh", host, "echo", random)
		pr, err := sc.ExecuteFullySilent(l)
		if err != nil {
			return err
		}
		if !pr.Successful() {
			return fmt.Errorf("host %s is not reachable via ssh", host)
		}
		if !strings.Contains(pr.TrimmedOutput(), random) {
			return fmt.Errorf("connection to host %s behaved unexpectedly via ssh", host)
		}
	}
	return nil
}
