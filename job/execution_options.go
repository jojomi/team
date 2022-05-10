package job

type ExecutionOptions struct {
	ExecutionPlan ExecutionPlan
	Delay         bool
	Wait          bool

	DryRun  bool
	Verbose bool
}
