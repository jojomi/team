package main

import (
	"fmt"
	"github.com/jojomi/team/job"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func getNextCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "next",
		Run:   handleNextCmd,
		Short: "Show jobs ordered by next execution",
	}

	// command local flags
	f := cmd.Flags()
	addJobSourceFlags(f)
	f.BoolP("unattended-only", "u", false, "only show jobs that run unattendedly")
	f.BoolP("possible-only", "p", false, "only show jobs that are currently possible")
	f.BoolP("fixable", "f", true, "only show jobs that are currently possible or fixable")
	f.Bool("show-unfixable", true, "show jobs that did not execute because of unfixable problems")

	return cmd
}

func handleNextCmd(cmd *cobra.Command, args []string) {
	var err error
	env := EnvNext{}
	err = env.ParseFrom(cmd, args)
	if err != nil {
		log.Fatal().Err(err).Msg("could not parse command")
	}
	err = handleNext(env)
	if err != nil {
		panic(err)
	}
}

func handleNext(env EnvNext) error {
	setLoggerVerbosity(env.Verbose)

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

	// sort
	pool = pool.SortedBy(func(j1, j2 job.Job) bool {
		t1, err := getJobNextDate(j1)
		if err != nil {
			return true
		}
		t2, err := getJobNextDate(j2)
		if err != nil {
			return false
		}
		return t1.Before(t2)
	})

	// print jobs in order
	for _, j := range pool.Jobs() {
		nextTime, err := getJobNextDate(j)
		if err != nil {
			log.Fatal().Err(err).Msgf("could not get next date for job %s", j.Metadata().Name)
		}
		job.PrintHeader(j, &nextTime)
		fmt.Println()
		fmt.Println()
	}

	return nil
}
