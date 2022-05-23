package main

import (
	"github.com/rs/zerolog/log"
)

func main() {
	setupLogger()

	// build root command
	rootCmd := getRootCmd()

	// add version command
	rootCmd.AddCommand(getVersionCmd())

	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("")
	}
}
