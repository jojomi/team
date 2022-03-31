package assert

import (
	"fmt"
	"github.com/jojomi/go-script/v2"
	"github.com/jojomi/team/exit"
	"math/rand"
	"strings"
)

func SSHReachableNonInteractive(args []string) (error, exit.Code) {
	var (
		sshOpts = map[string]string{
			"BatchMode":                       "yes",
			"ConnectTimeout":                  "5",
			"PubkeyAuthentication":            "yes",
			"PasswordAuthentication":          "no",
			"KbdInteractiveAuthentication":    "no",
			"ChallengeResponseAuthentication": "no",
		}
		sc     = script.NewContext()
		l      *script.LocalCommand
		random string
	)
	for _, host := range args {
		random = randSeq(32)
		l = sshCommandWithOpts(host, sshOpts)
		l.AddAll("echo", random)

		pr, err := sc.ExecuteFullySilent(l)
		if err != nil {
			return err, exit.CodeErrorFixable
		}
		if !pr.Successful() {
			return fmt.Errorf("host %s is not reachable non-interactively via ssh", host), exit.CodeErrorFixable
		}
		if !strings.Contains(pr.TrimmedOutput(), random) {
			return fmt.Errorf("connection to host %s behaved unexpectedly via ssh", host), exit.CodeErrorFinal
		}
	}

	return nil, exit.CodeOkay
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
