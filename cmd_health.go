package main

import (
	"fmt"
	"github.com/jojomi/team/job"
	"github.com/jojomi/team/persistance"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func getHealthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "health",
		Run: handleHealthCmd,
	}

	f := cmd.Flags()
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

	for _, j := range pool.Jobs() {
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
