package job

type ExecutionOptions struct {
	SkipExecution bool
	Delay         bool
	Wait          bool

	DryRun  bool
	Verbose bool
}
