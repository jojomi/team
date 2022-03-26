package main

import (
	"fmt"
	"github.com/jojomi/go-script/v2/print"
	"github.com/jojomi/team/job"
	"github.com/jojomi/team/persistance"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"time"
)

func getHealthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "health",
		Run: handleHealthCmd,
	}

	f := cmd.Flags()
	f.BoolP("required-only", "r", false, "only include currently required jobs in output")
	f.BoolP("possible-only", "p", false, "only include currently possible jobs in output")
	addJobSourceFlags(f)

	return cmd
}

func handleHealthCmd(cmd *cobra.Command, args []string) {
	env := EnvHealth{}
	err := env.ParseFrom(cmd, args)
	if err != nil {
		log.Fatal().Err(err).Msg("could not parse command")
	}
	err = handleHealth(env)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func handleHealth(env EnvHealth) error {
	setLoggerVerbosity(env.Verbose)

	pool := getCLIPool(env.JobDir, env.JobFile)

	if env.RequiredOnly {
		pool = pool.RequiredOnly(func(j job.Job) time.Time {
			t, err := getJobNextDate(j)
			if err != nil {
				log.Fatal().Err(err).Msgf("could not get next date for job %s", j.Metadata().Name)
			}
			return t
		})
	}
	if env.PossibleOnly {
		pool = pool.PossibleOnly()
	}

	count := len(pool.Jobs())
	for i, j := range pool.Jobs() {
		print.Boldf("[%d/%d] ", i+1, count)
		job.PrintHeader(j, nil)

		fmt.Print("\nLetzter erfolgreicher Lauf: ")
		lastSuccessfulRun, err := persistance.GetLastSuccessfulJobRun(j)
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}
		if lastSuccessfulRun == nil {
			fmt.Println("NIE")
		} else {
			fmt.Printf("%s (%s)\n", lastSuccessfulRun.Start.Format(job.DateFormat), job.Timeago(lastSuccessfulRun.Start))
		}

		lastRun, err := persistance.GetLastJobRun(j)
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}

		if lastSuccessfulRun == nil || lastRun == nil || lastRun.Start.Format(job.DateFormat) != lastSuccessfulRun.Start.Format(job.DateFormat) {
			fmt.Print("Letzter Lauf: ")
			if lastRun == nil {
				fmt.Println("NIE")
			} else {
				fmt.Printf("%s (%s)\n", lastRun.Start.Format(job.DateFormat), job.Timeago(lastRun.Start))
			}
		}

		nextDate, err := getJobNextDate(j)
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}
		fmt.Print("Nächster Lauf frühestens: ")
		fmt.Printf("%s (%s)\n", nextDate.Format(job.DateFormat), job.Timeago(nextDate))

		fmt.Println()
	}

	return nil
}
