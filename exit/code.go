package exit

type Code int

const (
	CodeOkay         Code = 0
	CodeErrorFixable      = iota + 1
	CodeErrorFinal
)
