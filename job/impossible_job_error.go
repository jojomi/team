package job

import (
	"fmt"
	"github.com/jojomi/go-script/v2"
	"github.com/jojomi/team/exit"
	"github.com/juju/errors"
	"strings"
)

type ImpossibleJobError struct {
	base    error
	fixable bool
}

func NewImpossibleJobError(err error, fixable bool) ImpossibleJobError {
	return ImpossibleJobError{
		base:    err,
		fixable: fixable,
	}
}

func NewImpossibleJobErrorFromProcessResult(result *script.ProcessResult) ImpossibleJobError {
	text := strings.TrimSpace(result.Error())
	exitCode, _ := result.ExitCode()
	fixable := exitCode == exit.CodeErrorFixable

	return NewImpossibleJobError(errors.New(text), fixable)
}

func (x ImpossibleJobError) IsFixable() bool {
	return x.fixable
}

func (x ImpossibleJobError) Error() string {
	if x.base == nil {
		return "no more details available"
	}
	return fmt.Sprintf("%s", x.base.Error())
}

func (x ImpossibleJobError) Unwrap() error {
	return x.base // Returns inner error
}

func (x *ImpossibleJobError) Is(target error) bool {
	_, ok := target.(*ImpossibleJobError)
	return ok
}

func (x *ImpossibleJobError) As(t any) bool {
	fmt.Println("checking as")
	_, ok := t.(*ImpossibleJobError)
	if !ok {
		return false
	}
	t = x
	return true
}
