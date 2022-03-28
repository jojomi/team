package main

import (
	"errors"
	"fmt"
	"github.com/jojomi/team/exit"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/jojomi/go-script/v2"

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
		"command-exists": assertCommandExists,
		"file-exists":    assertFileExists,
		"dir-exists":     assertDirExists,
		"mounted":        assertMounted,
		"non-empty-dir":  assertNonEmptyDir,
		"ssh-reachable":  assertSSHReachable,
		"weekday":        assertWeekday,
	}

	cmd := env.Command
	f, ok := handlerMap[cmd]
	if !ok {
		return fmt.Errorf("assert handler not found: %s", cmd), exit.CodeErrorFinal
	}

	return f(env.Args)
}

func assertCommandExists(args []string) (error, exit.Code) {
	if len(args) < 1 {
		return errors.New("no command name given"), exit.CodeErrorFinal
	}
	sc := script.NewContext()
	for _, cmd := range args {
		if sc.CommandExists(cmd) {
			continue
		}
		return fmt.Errorf("%s not installed or not in PATH", cmd), exit.CodeErrorFixable
	}
	return nil, exit.CodeOkay
}

func assertFileExists(args []string) (error, exit.Code) {
	if len(args) < 1 {
		return errors.New("no filename name given"), exit.CodeErrorFinal
	}
	sc := script.NewContext()
	for _, filename := range args {
		if sc.FileExists(filename) {
			continue
		}
		return fmt.Errorf("file %s does not exist", filename), exit.CodeErrorFixable
	}
	return nil, exit.CodeOkay
}

func assertMounted(args []string) (error, exit.Code) {
	if len(args) < 1 {
		return errors.New("no mountpoint dir given"), exit.CodeErrorFinal
	}
	sc := script.NewContext()
	for _, filename := range args {
		if !sc.DirExists(filename) {
			return fmt.Errorf("mounting dir %s does not exist", filename), exit.CodeErrorFixable
		}

	}
	return nil, exit.CodeOkay
}

func assertDirExists(args []string) (error, exit.Code) {
	if len(args) < 1 {
		return errors.New("no dir(s) given"), exit.CodeErrorFinal
	}
	sc := script.NewContext()

	for _, dir := range args {
		if !sc.DirExists(dir) {
			return fmt.Errorf("directory %s does not exist", dir), exit.CodeErrorFixable
		}
	}
	return nil, exit.CodeOkay
}

func assertNonEmptyDir(args []string) (error, exit.Code) {
	if len(args) < 1 {
		return errors.New("no dir(s) given"), exit.CodeErrorFinal
	}
	for _, dir := range args {
		if err, exitCode := assertDirExists([]string{dir}); err != nil {
			return err, exitCode
		}

		f, err := os.Open(dir)
		if err != nil {
			return err, exit.CodeErrorFixable
		}

		_, err = f.Readdirnames(1)
		if err == io.EOF {
			return fmt.Errorf("directory %s is empty", dir), exit.CodeErrorFixable
		}

		err = f.Close()
		if err != nil {
			return err, exit.CodeErrorFixable
		}
	}
	return nil, exit.CodeOkay
}

func assertSSHReachable(args []string) (error, exit.Code) {
	sc := script.NewContext()
	var (
		l      *script.LocalCommand
		random = randSeq(32)
	)
	for _, host := range args {
		l = script.NewLocalCommand()
		l.AddAll("ssh", host, "echo", random)
		pr, err := sc.ExecuteFullySilent(l)
		if err != nil {
			return err, exit.CodeErrorFixable
		}
		if !pr.Successful() {
			return fmt.Errorf("host %s is not reachable via ssh", host), exit.CodeErrorFixable
		}
		if !strings.Contains(pr.TrimmedOutput(), random) {
			return fmt.Errorf("connection to host %s behaved unexpectedly via ssh", host), exit.CodeErrorFinal
		}
	}
	return nil, exit.CodeOkay
}

func assertWeekday(args []string) (error, exit.Code) {
	weekday := time.Now().Weekday()
	dayString := strings.ToLower(weekday.String())

	// TODO parse -

	// one of?
	for _, arg := range args {
		if strings.EqualFold(arg, dayString) {
			return nil, exit.CodeOkay
		}
	}
	return fmt.Errorf("today is not in %s", strings.Join(args, ", ")), exit.CodeErrorFinal
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
