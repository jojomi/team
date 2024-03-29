package job

import (
	"errors"
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
		possible, err := IsPossible(job)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Job %s: %s\n", job.Metadata().Name, err.Error())
			return false
		}
		return possible
	})
	return NewJobPoolWithJobs(filteredJobs)
}

func (x *Pool) PossibleOrFixableOnly() *Pool {
	filteredJobs := arrayMap[Job](x.jobs, func(job Job) bool {
		possible, err := IsPossible(job)

		if possible {
			return true
		}

		if e, ok := err.(ImpossibleJobError); ok {
			var jobErr ImpossibleJobError
			if errors.As(e, &jobErr) {
				if jobErr.IsFixable() {
					return true
				}
			}
			return false
		}
		return false
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

func (x *Pool) AttendedOnly() *Pool {
	filteredJobs := arrayMap[Job](x.jobs, func(job Job) bool {
		return !job.Metadata().IsUnattended()
	})
	return NewJobPoolWithJobs(filteredJobs)
}

func (x *Pool) EnabledOnly() *Pool {
	filteredJobs := arrayMap[Job](x.jobs, func(job Job) bool {
		return job.Metadata().IsEnabled()
	})
	return NewJobPoolWithJobs(filteredJobs)
}

func (x *Pool) SortedBy(lessFunc func(j1, j2 Job) bool) *Pool {
	sortedJobs := x.jobs
	sort.SliceStable(sortedJobs, func(i, j int) bool {
		return lessFunc(x.jobs[i], x.jobs[j])
	})
	return NewJobPoolWithJobs(sortedJobs)
}

func (x *Pool) ByWeight() *Pool {
	return x.SortedBy(func(j1, j2 Job) bool {
		return j1.Metadata().GetWeight() > j2.Metadata().GetWeight()
	})
}

// NeedsSudo checks if any of the Jobs in the Pool needs privileged access
func (x Pool) NeedsSudo() bool {
	for _, job := range x.jobs {
		if s := job.Metadata().Sudo; s != nil && *s == true {
			return true
		}
	}
	return false
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
