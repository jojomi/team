package assert

import (
	"fmt"
	"github.com/jojomi/go-script/v2"
	"github.com/jojomi/team/exit"
	"github.com/juju/errors"
)

func CommandExists(args []string) (error, exit.Code) {
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
