package assert

import (
	"fmt"
	"github.com/jojomi/go-script/v2"
	"github.com/jojomi/team/exit"
	"strings"
)

func SSHReachable(args []string) (error, exit.Code) {
	var (
		sshOpts = map[string]string{
			"BatchMode":                       "yes",
			"ConnectTimeout":                  "5",
			"PubkeyAuthentication":            "no",
			"PasswordAuthentication":          "no",
			"KbdInteractiveAuthentication":    "no",
			"ChallengeResponseAuthentication": "no",
		}
		sc = script.NewContext()
		l  *script.LocalCommand
	)

	for _, host := range args {
		l = sshCommandWithOpts(host, sshOpts)

		pr, err := sc.ExecuteFullySilent(l)
		if err != nil {
			return err, exit.CodeErrorFixable
		}
		if !pr.Successful() {
			// the "good case" ("permission denied" because we did not even try to authenticate)
			if strings.Contains(pr.Error(), "Permission denied") {
				continue
			}
			return fmt.Errorf("host %s is not reachable via ssh", host), exit.CodeErrorFixable
		}
	}
	return nil, exit.CodeOkay
}
