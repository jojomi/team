package assert

import (
	"fmt"
	"github.com/jojomi/team/exit"
	"strings"
	"time"
)

func TimeBefore(args []string) (error, exit.Code) {
	n := time.Now().Format(timeFormat)

	// one of?
	for _, arg := range args {
		if !isValidTime(arg) {
			return fmt.Errorf("invalid time format: %s, format is %s", arg, timeFormat), exit.CodeErrorFinal
		}
		if n <= arg {
			return nil, exit.CodeOkay
		}
	}
	return fmt.Errorf("time is not before %s", strings.Join(args, ", ")), exit.CodeErrorFinal
}
