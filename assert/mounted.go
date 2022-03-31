package assert

import (
	"fmt"
	"github.com/jojomi/go-script/v2"
	"github.com/jojomi/team/exit"
	"github.com/juju/errors"
)

func Mounted(args []string) (error, exit.Code) {
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
