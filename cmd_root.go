package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"time"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/jojomi/team/persistance"

	"github.com/PaesslerAG/gval"
	"github.com/google/uuid"
	"github.com/iancoleman/strcase"
	"github.com/jojomi/go-script/v2"
	"github.com/jojomi/go-script/v2/interview"
	"github.com/jojomi/go-script/v2/print"
	"github.com/jojomi/team/ent/run"
	"github.com/jojomi/team/job"
	jujuErrors "github.com/juju/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func getRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: strcase.ToKebab(ToolName),
		Run: handleRootCmd,
	}

	cmd.AddCommand(getNextCmd(), getHealthCmd(), getVersionCmd())

	pf := cmd.PersistentFlags()
	pf.BoolP("verbose", "v", false, "activate verbose output")

	// command local flags
	f := cmd.Flags()
	addJobSourceFlags(f)
	f.BoolP("manual", "m", true, "go job by job manually")
	f.BoolP("select-first", "s", false, "first select, then execute jobs")
	f.BoolP("default-yes", "y", false, "set default answer to yes")
	f.BoolP("unattended-only", "u", false, "only offer jobs that run unattendedly")
	f.BoolP("possible-only", "p", false, "only offer jobs that are currently possible")
	f.BoolP("fixable", "f", true, "only offer jobs that are currently possible or fixable")
	f.BoolP("skip-time-check", "t", false, "don't check if the job is required by time")
	f.BoolP("attended-first", "a", false, "optimize to do attended tasks as early as possible") // TODO implement
	f.Bool("show-unfixable", false, "show jobs that did not execute because of unfixable problems")
	f.BoolP("dry-run", "d", false, "prevent destructive operations")
	f.Bool("shutdown", false, "shutdown after completing (with a timeout, so you can stop it)")

	return cmd
}

func handleRootCmd(cmd *cobra.Command, args []string) {
	var err error
	env := EnvRoot{}
	err = env.ParseFrom(cmd, args)
	if err != nil {
		log.Fatal().Err(err).Msg("could not parse command")
	}
	err = handleRoot(env)
	if err != nil {
		panic(err)
	}
}

func handleRoot(env EnvRoot) error {
	setLoggerVerbosity(env.Verbose)

	var (
		executionOptionMap = make(map[job.Job]job.ExecutionOptions, 0)
	)

	totalStart := time.Now()
	handledJobs := 0

	pool := getCLIPool(env.JobDir, env.JobFile)
	if env.UnattendedOnly {
		pool = pool.UnattendedOnly()
	}
	if env.PossibleOnly {
		pool = pool.PossibleOnly()
	}
	if !env.ShowUnfixable {
		pool = pool.PossibleOrFixableOnly()
	}

	// filter
	requiredPool := pool
	if !env.SkipTimeCheck {
		requiredPool = job.NewJobPool()
		now := time.Now()
		for _, j := range pool.Jobs() {
			next, err := getJobNextDate(j)
			if err != nil {
				log.Fatal().Err(err).Msg("")
			}
			if now.Before(next) {
				continue
			}

			requiredPool.AddJob(j)
		}
	}

	delayedJobs := job.NewJobPool()

	// loop relevant jobs in order
	poolSize := len(requiredPool.Jobs())
	for i, j := range requiredPool.Jobs() {
		executionOptions, err := handleJobPreparation(j, env, i, poolSize)
		if errors.Is(err, UserInterruptedError{}) {
			fmt.Println()
			break
		}
		if errors.Is(err, UserAbortedError{}) {
			fmt.Println()
			continue
		}
		if executionOptions.ExecutionPlan == job.ExecutionPlanSkip {
			fmt.Println()
			continue
		}
		if errors.Is(err, job.ImpossibleJobError{}) {
			job.PrintUnsuccessful("", nil)
			fmt.Println()
			continue
		}
		executionOptionMap[j] = executionOptions

		// execute now or later?
		if env.SelectFirst || executionOptions.Delay {
			delayedJobs.AddJob(j)
			fmt.Println()
			continue
		}

		start := time.Now()
		err = handleJobExecution(j, executionOptions, env)
		handledJobs++
		diff := time.Now().Sub(start).Round(time.Second)
		if err != nil {
			if _, ok := err.(job.ImpossibleJobError); !ok {
				log.Error().Err(err).Str("filename", j.Metadata().Name).Msg("could not handle job")
			}
			job.PrintUnsuccessful("", &diff)
			fmt.Println()
			continue
		}
		job.PrintSuccessful("", &diff)
		fmt.Println()
	}

	// execution now if it was delayed
	if len(delayedJobs.Jobs()) > 0 {
		var (
			count = len(delayedJobs.Jobs())
			err   error
		)

		fmt.Println("================================================")
		fmt.Println("               EXECUTION PHASE")
		fmt.Println("================================================")

		for i, j := range delayedJobs.Jobs() {
			print.Boldf("[%d/%d] ", i+1, count)
			job.PrintHeader(j, nil)
			fmt.Println()

			start := time.Now()

			executionOptions, ok := executionOptionMap[j]
			if !ok {
				executionOptions = job.ExecutionOptions{
					Wait:          true,
					ExecutionPlan: job.ExecutionPlanExecute,
					DryRun:        env.DryRun,
					Verbose:       env.Verbose,
				}
			}

			err = handleJobExecution(j, executionOptions, env)
			handledJobs++
			diff := time.Now().Sub(start).Round(time.Second)
			if errors.Is(err, UserInterruptedError{}) {
				fmt.Println()
				break
			}
			if errors.Is(err, UserAbortedError{}) {
				fmt.Println()
				continue
			}
			if err != nil {
				log.Error().Err(err).Str("filename", j.Metadata().Name).Msg("could not handle job")
				job.PrintUnsuccessful("", &diff)
				fmt.Println()
				continue
			}
			job.PrintSuccessful("", &diff)
			fmt.Println()
		}
	}

	totalDiff := time.Now().Sub(totalStart).Round(time.Second)
	fmt.Printf("Worked on %d tasks for %s\n", handledJobs, totalDiff)

	if env.Shutdown {
		shutdown(env)
	}

	return nil
}

func handleJobPreparation(j job.Job, env EnvRoot, index, count int) (job.ExecutionOptions, error) {
	opts := job.ExecutionOptions{
		ExecutionPlan: job.ExecutionPlanExecute,
		Delay:         false,
		Wait:          true,

		DryRun:  env.DryRun,
		Verbose: env.Verbose,
	}

	print.Boldf("[%d/%d] ", index+1, count)
	next, err := getJobNextDate(j)
	if err != nil {
		return opts, errors.New("could not get next execution date")
	}
	job.PrintHeader(j, &next)
	fmt.Println()

	if env.UnattendedOnly {
		possible, err := job.IsPossible(j)
		if errors.Is(err, job.ImpossibleJobError{}) {
			return opts, err
		}
		if err != nil {
			return opts, job.NewImpossibleJobError(err, false)
		}
		if !possible {
			return opts, job.NewImpossibleJobError(nil, false)
		}
	}

	var (
		possible bool
		jobError error
		again    bool
	)
	for {
		possible, jobError = job.IsPossible(j)

		if possible {
			break
		}

		print.Bold("Job not possible")
		fmt.Println(", reason:")

		if jobError != nil {
			fmt.Println(jobError)
		} else if !possible {
			fmt.Println("no further details")
		}

		again = false
		if e, ok := jobError.(job.ImpossibleJobError); ok {
			var jobErr job.ImpossibleJobError
			if errors.As(e, &jobErr) {
				if jobErr.IsFixable() {
					again, err = interview.Confirm("Try again?", false)
					if err != nil {
						return opts, err
					}
				}
			}
		}
		if !again {
			return opts, job.NewImpossibleJobError(jobError, false)
		}
	}

	if env.Manual {
		executeAction := &interview.Action{
			Label: "Execute",
		}
		doneAction := &interview.Action{
			Label: "Done",
		}
		executeAndDoneAction := &interview.Action{
			Label: "Execute and done",
		}
		doLaterAction := &interview.Action{
			Label: "Do Later",
		}
		skipAction := &interview.Action{
			Label: "Skip",
		}
		actions := []*interview.Action{
			executeAction,
			doneAction,
			executeAndDoneAction,
			doLaterAction,
			skipAction,
		}

		defaultAction := skipAction
		if env.DefaultYes {
			defaultAction = executeAction
		}

		actions = interview.WithAutoShortcuts(actions)
		action, err := interview.SelectActionWithDefault(actions, defaultAction)
		// doExec, err = interview.Confirm("Aufgabe ausführen?", env.DefaultYes)
		if errors.Is(err, terminal.InterruptErr) {
			return opts, UserInterruptedError{}
		}
		if err != nil {
			log.Error().Err(err).Msg("")
		}

		switch action {
		case skipAction:
			opts.ExecutionPlan = job.ExecutionPlanSkip
		case executeAndDoneAction:
			opts.Wait = false
			opts.ExecutionPlan = job.ExecutionPlanExecute
		case doneAction:
			opts.Wait = false
			opts.ExecutionPlan = job.ExecutionPlanLogDone
		case doLaterAction:
			opts.Delay = true
		}
	}

	return opts, nil
}

func handleJobExecution(j job.Job, executionOptions job.ExecutionOptions, env EnvRoot) error {
	var (
		err error
		id  uuid.UUID
	)

	// is this job possible?
	possible, err := job.IsPossible(j)
	if e, ok := err.(job.ImpossibleJobError); ok {
		return e
	}
	if err != nil {
		return job.NewImpossibleJobError(err, false)
	}
	if !possible {
		return job.NewImpossibleJobError(nil, false)
	}

	if !env.DryRun {
		id, err = logJobExecutionStart(j)
		if err != nil {
			return err
		}
	}

	if executionOptions.ExecutionPlan == job.ExecutionPlanExecute {
		err = j.Pre()
		if err != nil {
			// TODO log
			return jujuErrors.Annotate(err, "could not pre")
		}

		err = j.Execute(executionOptions)

		err = j.Post()
		if err != nil {
			// TODO log
			return jujuErrors.Annotate(err, "could not post")
		}
	}

	if !env.DryRun {
		logError := logJobExecution(id, j, err, "" /* TODO */)
		if logError != nil {
			return logError
		}
	}

	if err != nil {
		return jujuErrors.Annotate(err, "could not execute")
	}

	return nil
}

func logJobExecutionStart(j job.Job) (uuid.UUID, error) {
	db := persistance.GetDatabaseClient()
	x := db.Run.Create()

	uuidValue := uuid.New()

	x.SetID(uuidValue)
	x.SetJob(j.Metadata().Name)
	x.SetStart(time.Now())
	x.SetStatus(run.StatusRunning)

	err := x.Exec(context.Background())
	return uuidValue, err
}

func logJobExecution(id uuid.UUID, j job.Job, err error, output string) error {
	db := persistance.GetDatabaseClient()
	x := db.Run.UpdateOneID(id)

	x.SetEnd(time.Now()).SetLog(output)

	if err == nil {
		x.SetStatus(run.StatusSuccessful)
	} else {
		x.SetStatus(run.StatusFailed)
	}

	_, err = x.Save(context.Background())
	return err
}

func shutdown(env EnvRoot) {
	const delayMinutes = 5

	fmt.Println()
	fmt.Println()
	if env.DryRun {
		job.PrintDryRunTag()
	}
	print.Boldf("Shutdown in %d minutes... (abort with Ctrl+C)\n", delayMinutes)

	// catch ctrl+c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("\nShutdown stopped.")
		os.Exit(1)
	}()

	time.Sleep(delayMinutes * time.Minute)

	if env.DryRun {
		print.Bold("[DRY-RUN] ")
	}
	fmt.Println("Shutting down now.")
	if env.DryRun {
		return
	}

	sc := script.NewContext()
	cmd := script.LocalCommandFrom("systemctl poweroff")
	_, err := sc.ExecuteDebug(cmd)
	if err != nil {
		log.Fatal().Err(err).Msg("could not shutdown")
	}
}

func getJobNextDate(j job.Job) (time.Time, error) {
	next := j.Metadata().Next
	if next == "" {
		return time.Now(), nil
	}

	var (
		lastTime        time.Time
		lastSuccessTime time.Time
	)
	lastSuccess, err := persistance.GetLastSuccessfulJobRun(j)
	if err != nil {
		return time.Time{}, fmt.Errorf("could not get last successful run for job %s", j.Metadata().Name)
	}
	if lastSuccess == nil {
		lastSuccessTime = time.Time{}
	} else {
		lastSuccessTime = lastSuccess.Start
	}

	last, err := persistance.GetLastJobRun(j)
	if err != nil {
		return time.Time{}, fmt.Errorf("could not get last run for job %s", j.Metadata().Name)
	}
	if last == nil {
		lastTime = time.Time{}
	} else {
		lastTime = last.Start
	}

	variables := map[string]interface{}{
		"lastSuccess": lastSuccessTime,
		"last":        lastTime,
	}

	dateAdder := gval.InfixOperator("+", func(a, b interface{}) (interface{}, error) {
		date, ok1 := a.(time.Time)
		if !ok1 {
			return nil, fmt.Errorf("unexpected syntax in %s: casting to time failed", next)
		}

		durationString, ok2 := b.(string)
		if !ok2 {
			return nil, fmt.Errorf("unexpected syntax in %s: casting to string failed", next)
		}
		duration, err := parseDuration(durationString)
		if err != nil {
			return nil, fmt.Errorf("invalid duration (%s): %s", durationString, err.Error())
		}
		return date.Add(duration), nil
	})

	value, err := gval.Evaluate(next,
		variables,
		dateAdder,
	)
	if err != nil {
		return time.Time{}, err
	}

	return value.(time.Time), nil
}

func parseDuration(input string) (time.Duration, error) {
	re := regexp.MustCompile(`(?:(\d+)(y))?(?:(\d+)(w))?(?:(\d+)(d))?(?:(\d+)(h))?(?:(\d+)(m))?`)
	parts := re.FindStringSubmatch(input)
	if parts == nil {
		return time.Millisecond, fmt.Errorf("invalid duration specification: %s", input)
	}
	seconds := 0
	factors := map[string]int{
		"m": 60,
		"h": 60 * 60,
		"d": 60 * 60 * 24,
		"y": 60 * 60 * 24 * 365.25,
	}

	for i := 2; i < len(parts); i = i + 2 {
		einheit := parts[i]
		if einheit == "" {
			continue
		}
		factor, ok := factors[einheit]
		if !ok {
			return time.Millisecond, fmt.Errorf("invalid duration specification (factor not found for %s): %s", einheit, input)
		}
		value, err := strconv.Atoi(parts[i-1])
		if err != nil {
			return time.Millisecond, fmt.Errorf("invalid duration specification: %s", input)
		}
		seconds += factor * value
	}

	result, err := time.ParseDuration(fmt.Sprintf("%ds", seconds))
	if err != nil {
		return time.Millisecond, fmt.Errorf("invalid duration specification: %s", input)
	}
	return result, nil
}
