package main

import (
	"github.com/rs/zerolog/log"
)

func main() {
	setupLogger()

	// build root command
	rootCmd := getRootCmd()

	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("")
	}
}
