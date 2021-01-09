package retry

// Error constant texts
const (
	CannotBeRetried    = "this operation cannot be retried"
	MaxRetriesExceeded = "this operation exceeded the maximum number of retries"
	MaxElapsedExceeded = "this operation exceeded the maximum time allowed to complete"
)

// ErrDoNotRetry is used to wrap errors from retried calls that shouldn't be retried
type ErrDoNotRetry struct {
	err error
}

// DoNotRetry wraps an error with ErrDoNotRetry so that it won't be retried by a Retryer
func DoNotRetry(err error) error {
	return &ErrDoNotRetry{err: err}
}

func (e *ErrDoNotRetry) Error() string {
	return e.err.Error()
}

func (e *ErrDoNotRetry) Unwrap() error {
	return e.err
}
