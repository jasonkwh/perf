package perf

import "fmt"

var (
	ErrFailedServerListenerInit = func(err error) error {
		return fmt.Errorf("failed to initialise server listener: %w", err)
	}
)
