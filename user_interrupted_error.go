package main

type UserInterruptedError struct {
}

func (x UserInterruptedError) Error() string {
	return "user interrupted"
}
