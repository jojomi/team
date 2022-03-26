package job

import (
	"fmt"
	"os"
	"sort"
	"time"
)

type Pool struct {
	jobs []Job
}

func NewJobPool() *Pool {
	return NewJobPoolWithJobs([]Job{})
}

func NewJobPoolWithJobs(jobs []Job) *Pool {
	return &Pool{
		jobs: jobs,
	}
}

func (x *Pool) AddJob(job Job) *Pool {
	x.jobs = append(x.jobs, job)
	return x
}

func (x *Pool) Jobs() []Job {
	return x.jobs
}

func (x *Pool) PossibleOnly() *Pool {
	filteredJobs := arrayMap[Job](x.jobs, func(job Job) bool {
		possible, err := job.IsPossible()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Job %s: %s\n", job.Metadata().Name, err.Error())
			return false
		}
		return possible
	})
	return NewJobPoolWithJobs(filteredJobs)
}

func (x *Pool) RequiredOnly(next func(j Job) time.Time) *Pool {
	now := time.Now()
	filteredJobs := arrayMap[Job](x.jobs, func(job Job) bool {
		nextTime := next(job)
		return !now.Before(nextTime)
	})
	return NewJobPoolWithJobs(filteredJobs)
}

func (x *Pool) UnattendedOnly() *Pool {
	filteredJobs := arrayMap[Job](x.jobs, func(job Job) bool {
		return job.Metadata().IsUnattended()
	})
	return NewJobPoolWithJobs(filteredJobs)
}

func (x *Pool) EnabledOnly() *Pool {
	filteredJobs := arrayMap[Job](x.jobs, func(job Job) bool {
		return job.Metadata().IsEnabled()
	})
	return NewJobPoolWithJobs(filteredJobs)
}

func (x *Pool) ByWeight() *Pool {
	sortedJobs := x.jobs
	sort.SliceStable(sortedJobs, func(i, j int) bool {
		return x.jobs[i].Metadata().GetWeight() > x.jobs[j].Metadata().GetWeight()
	})
	return NewJobPoolWithJobs(sortedJobs)
}

func arrayMap[T any](input []T, keep func(item T) bool) []T {
	result := make([]T, 0)

	for _, elem := range input {
		if !keep(elem) {
			continue
		}
		result = append(result, elem)
	}

	return result
}
