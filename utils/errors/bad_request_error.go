package errors

type BadRequestError struct {
	message string
}

func NewBadRequestError(message string) BadRequestError {
	return BadRequestError{
		message: message,
	}
}

func (e BadRequestError) Error() string {
	return e.message
}
