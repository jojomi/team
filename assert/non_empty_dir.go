package assert

import (
	"fmt"
	"github.com/jojomi/team/exit"
	"github.com/juju/errors"
	"io"
	"os"
)

func NonEmptyDir(args []string) (error, exit.Code) {
	if len(args) < 1 {
		return errors.New("no dir(s) given"), exit.CodeErrorFinal
	}
	for _, dir := range args {
		if err, exitCode := DirExists([]string{dir}); err != nil {
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
