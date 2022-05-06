package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jojomi/team/job"
	jujuErrors "github.com/juju/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
)

func addJobSourceFlags(f *pflag.FlagSet) {
	f.String("job-dir", "~/.team/jobs", "directory with the job files")
	f.StringP("job-file", "j", "", "handle only one specific job file")
}

func getCLIPool(jobDir, jobFile string) *job.Pool {
	var (
		jobFiles []string
		err      error
	)
	if jobFile != "" {
		var f string
		if filepath.IsAbs(jobFile) {
			f = jobFile
		} else {
			f = filepath.Join(jobDir, jobFile)
		}
		if !strings.HasSuffix(strings.ToLower(f), ".yml") && !strings.HasSuffix(strings.ToLower(f), ".yaml") {
			f = f + ".yml"
		}
		jobFiles = []string{f}
	} else {
		jobFiles, err = getAllJobYAMLs(jobDir)
		if err != nil {
			log.Fatal().Err(err).Msg("could get job definitions")
		}
	}

	// fill job pool
	jobPool := job.NewJobPool()
	for _, filename := range jobFiles {
		factory := job.Factory{}
		j, err := factory.GetFromFile(filename)
		if err != nil {
			err = jujuErrors.Annotatef(err, "could not load job from %s", filename)
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			continue
		}
		jobPool.AddJob(j)
	}

	// filter & sort
	pool := jobPool.EnabledOnly()
	pool = pool.ByWeight()
	return pool
}

func getAllJobYAMLs(jobDir string) ([]string, error) {
	searchGlob := filepath.Join(jobDir, "*.yml")
	files, err := filepath.Glob(searchGlob)
	if err != nil {
		return []string{}, jujuErrors.Annotatef(err, "could not get job definitions with glob %s", searchGlob)
	}
	return files, nil
}
