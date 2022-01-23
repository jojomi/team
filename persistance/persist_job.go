package persistance

import (
	"context"

	"github.com/jojomi/team/ent"
	"github.com/jojomi/team/ent/run"
	"github.com/jojomi/team/job"
)

func GetLastJobRun(j job.Job) (*ent.Run, error) {
	client := GetDatabaseClient()
	run, err := client.Run.Query().Where(run.JobEQ(j.Metadata().Name)).Order(ent.Desc(run.FieldStart)).First(context.Background())
	if err != nil && err.Error() == "ent: run not found" {
		return nil, nil
	}
	/* does not work: if errors.Is(err, &ent.NotFoundError{}) {
		return nil, nil
	}*/
	return run, err
}

func GetLastSuccessfulJobRun(j job.Job) (*ent.Run, error) {
	client := GetDatabaseClient()
	run, err := client.Run.Query().Where(run.And(run.JobEQ(j.Metadata().Name), run.StatusEQ(run.StatusSuccessful))).Order(ent.Desc(run.FieldStart)).First(context.Background())
	if err != nil && err.Error() == "ent: run not found" {
		return nil, nil
	}
	/* does not work: if errors.Is(err, &ent.NotFoundError{}) {
		return nil, nil
	}*/
	return run, err
}
