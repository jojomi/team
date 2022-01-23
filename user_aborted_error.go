package main

type UserAbortedError struct {
}

func (x UserAbortedError) Error() string {
	return "user aborted"
}
