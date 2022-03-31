package assert

import (
	"fmt"
	"github.com/jojomi/go-script/v2"
	"github.com/jojomi/team/exit"
	"math/rand"
	"strings"
)

func SSHReachable(args []string) (error, exit.Code) {
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

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
