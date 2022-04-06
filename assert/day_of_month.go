package assert

import (
	"fmt"
	"github.com/jinzhu/now"
	"github.com/jojomi/team/exit"
	"strconv"
	"strings"
	"time"
)

func DayOfMonth(args []string) (error, exit.Code) {
	timeNow := time.Now()
	dayOfMonth := timeNow.Day()

	// one of?
	for _, arg := range args {
		d, err := strconv.Atoi(arg)
		if err != nil {
			return err, exit.CodeErrorFinal
		}

		// from end? 0 = last, -1 = last but one and so on
		if d < 1 {
			d = now.With(timeNow).EndOfMonth().AddDate(0, 0, d).Day()
		}

		if d == dayOfMonth {
			return nil, exit.CodeOkay
		}
	}
	return fmt.Errorf("today's day in month is not in %s", strings.Join(args, ", ")), exit.CodeErrorFinal
}
