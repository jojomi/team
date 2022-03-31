package assert

import (
	"fmt"
	"github.com/jojomi/go-script/v2"
	"github.com/jojomi/team/exit"
	"github.com/juju/errors"
)

func DirExists(args []string) (error, exit.Code) {
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
