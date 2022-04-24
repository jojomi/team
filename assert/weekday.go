package assert

import (
	"fmt"
	"github.com/jojomi/team/exit"
	"github.com/jojomi/team/ranges"
	"strings"
	"time"
)

func Weekday(args []string) (error, exit.Code) {
	weekday := time.Now().Weekday()
	dayStringLong := strings.ToLower(weekday.String())
	dayStringShort := dayStringLong[0:3]

	checks := []struct {
		current string
		values  []string
	}{
		{
			current: dayStringLong,
			values:  weekdayListLong,
		},
		{
			current: dayStringShort,
			values:  weekdayListShort,
		},
	}

	for _, check := range checks {
		resolvedArgs, err := ranges.ResolveIndexedRanges(args, check.values)
		if err != nil {
			return err, exit.CodeErrorFinal
		}

		// one of?
		for _, arg := range resolvedArgs {
			if strings.EqualFold(arg, check.current) {
				return nil, exit.CodeOkay
			}
		}
	}

	return fmt.Errorf("today is not in %s", strings.Join(args, ", ")), exit.CodeErrorFinal
}

var (
	weekdayListLong = []string{
		"monday",
		"tuesday",
		"wednesday",
		"thursday",
		"friday",
		"saturday",
		"sunday",
	}
	weekdayListShort = []string{
		"mon",
		"tue",
		"wed",
		"thu",
		"fri",
		"sat",
		"sun",
	}
)
