package assert

import (
	"fmt"
	"github.com/jojomi/team/exit"
	"strings"
	"time"
)

func Weekday(args []string) (error, exit.Code) {
	weekday := time.Now().Weekday()
	dayString := strings.ToLower(weekday.String())

	// TODO parse -

	// one of?
	for _, arg := range args {
		if strings.EqualFold(arg, dayString) {
			return nil, exit.CodeOkay
		}
	}
	return fmt.Errorf("today is not in %s", strings.Join(args, ", ")), exit.CodeErrorFinal
}
