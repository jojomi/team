package job

import (
	"fmt"
	"github.com/jojomi/go-script/v2"
	"github.com/jojomi/go-script/v2/print"
	"github.com/muesli/termenv"
	"github.com/xeonx/timeago"
	"net/http"
	"time"
)

func getExecutorByOutputType(output string) func(command script.Command) (pr *script.ProcessResult, err error) {
	sc := script.NewContext()
	if output == "debug" || output == "" {
		return sc.ExecuteDebug
	}
	return sc.ExecuteSilent
}

func isURLReachable(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		return false
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return false
	}
	return true
}

func PrintDryRunTag() {
	fmt.Print(termenv.String("[DRY-RUN] ").Foreground(cp.Color("#f5993d")))
}

func debugCommand(cmd script.Command, options ExecutionOptions) {
	if !options.Verbose && !options.DryRun {
		return
	}

	if options.DryRun {
		PrintDryRunTag()
	}

	fmt.Println("»", termenv.String(cmd.String()).Foreground(cp.Color("#aaaa00")))
}

func PrintSuccessful(message string, duration *time.Duration) {
	if message == "" {
		message = "erfolgreich"
	}

	c := cp.Color("#00cc66")

	fmt.Print(termenv.String("✓ " + message).Foreground(c))
	if duration != nil {
		fmt.Print(" – " + duration.String())
	}
	fmt.Println()
}

func PrintUnsuccessful(message string, duration *time.Duration) {
	if message == "" {
		message = "fehlgeschlagen"
	}

	c := cp.Color("#db4646")

	fmt.Print(termenv.String("× " + message).Foreground(c))
	if duration != nil {
		fmt.Print(" – " + duration.String())
	}
	fmt.Println()
}

func PrintHeader(j Job, next *time.Time) {
	print.Bold(termenv.String(j.Metadata().Name).Foreground(cp.Color("#809fff")))

	filename := j.Metadata().Filename
	if filename != "" {
		fmt.Printf(" (from %s)", filename)
	}

	if next != nil {
		fmt.Printf("\nFällig seit %s (%s)", next.Format(DateFormat), Timeago(*next))
	}
}

func Timeago(date time.Time) string {
	return timeago.German.Format(date)
}
