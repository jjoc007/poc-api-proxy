package errors

type NotFoundError struct {
	message string
}

func NewNotFoundError(message string) NotFoundError {
	return NotFoundError{
		message: message,
	}
}

func (e NotFoundError) Error() string {
	return e.message
}
