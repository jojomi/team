package main

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func getVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "version",
		Run: handleVersionCmd,
	}
	return cmd
}

func handleVersionCmd(cmd *cobra.Command, args []string) {
	version := GitVersion
	fmt.Println(ToolName + ", version " + version)

	v, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		log.Fatal().Err(err).Msg("flag parsing failed")
	}
	setLoggerVerbosity(v)

	log.Debug().
		Str("git version", GitVersion).
		Str("git commit", GitCommit).
		Str("git date", GitDate).
		Str("git state", GitState).
		Str("git branch", GitBranch).
		Str("git remote", GitRemote).
		Msg("")
}
