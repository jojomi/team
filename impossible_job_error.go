package main

import "fmt"

type ImpossibleJobError struct {
	base error
}

func NewImpossibleJobError(err error) ImpossibleJobError {
	return ImpossibleJobError{
		base: err,
	}
}

func (x ImpossibleJobError) Error() string {
	if x.base == nil {
		return "job is impossible (no more details available)"
	}
	return fmt.Sprintf("job is impossible: %s", x.base.Error())
}
