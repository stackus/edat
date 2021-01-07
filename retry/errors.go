package retry

const (
	CannotBeRetried    = "this operation cannot be retried"
	MaxRetriesExceeded = "this operation exceeded the maximum number of retries"
	MaxElapsedExceeded = "this operation exceeded the maximum time allowed to complete"
)

type ErrDoNotRetry struct {
	err error
}

func DoNotRetry(err error) error {
	return &ErrDoNotRetry{err: err}
}

func (e *ErrDoNotRetry) Error() string {
	return e.err.Error()
}

func (e *ErrDoNotRetry) Unwrap() error {
	return e.err
}
