package job

import "github.com/muesli/termenv"

var cp termenv.Profile

const DateFormat = "02.01.2006 15:04 Uhr"

func init() {
	cp = termenv.EnvColorProfile()
}
