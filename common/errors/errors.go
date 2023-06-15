package errors

//RequestError
// The following code defines a custom error type called RequestError with a status code and message.
// @property {int} StatusCode - StatusCode is an integer property that represents the HTTP status code
// of a request. It is typically a three-digit number that indicates the status of the request, such as
// 200 for a successful request or 404 for a not found error.
// @property {string} Message - The Message property is a string that represents the error message
// associated with the RequestError. It provides additional information about the error that occurred
// during the request.
type RequestError struct {
	StatusCode int
	Message    string
}

func (e RequestError) Error() string {
	return e.Message
}
