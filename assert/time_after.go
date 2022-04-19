package assert

import (
	"fmt"
	"github.com/jojomi/team/exit"
	"strings"
	"time"
)

const timeFormat = "15:04"

func TimeAfter(args []string) (error, exit.Code) {
	n := time.Now().Format(timeFormat)

	// one of?
	for _, arg := range args {
		if !isValidTime(arg) {
			return fmt.Errorf("invalid time format: %s, format is %s", arg, timeFormat), exit.CodeErrorFinal
		}
		if n >= arg {
			return nil, exit.CodeOkay
		}
	}
	return fmt.Errorf("time is not after %s", strings.Join(args, ", ")), exit.CodeErrorFinal
}

func isValidTime(arg string) bool {
	_, err := time.Parse(timeFormat, arg)
	return err == nil
}
