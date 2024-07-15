package exception

type NotFoundError struct {
	Error string
}

func NewNotFoundError(error string) NotFoundError {
	return NotFoundError{Error: error}
}

type DuplicateEmailError struct {
	Error string
}

func NewDuplicateEmail(error string) DuplicateEmailError {
	return DuplicateEmailError{Error: error}
}
