package errors

type FailedDependencyError struct {
	message string
}

func NewFailedDependencyError(message string) FailedDependencyError {
	return FailedDependencyError{
		message: message,
	}
}

func (e FailedDependencyError) Error() string {
	return e.message
}
