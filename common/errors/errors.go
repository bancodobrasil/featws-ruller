package errors

type RequestError struct {
	StatusCode int
	Message    string
}

func (e RequestError) Error() string {
	return e.Message
}
