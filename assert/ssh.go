package assert

import (
	"github.com/jojomi/go-script/v2"
)

func sshCommandWithOpts(host string, sshOpts map[string]string) *script.LocalCommand {
	l := script.LocalCommandFrom("ssh")
	for opt, optVal := range sshOpts {
		l.AddAll("-o", opt+"="+optVal)
	}
	l.Add(host)
	return l
}
