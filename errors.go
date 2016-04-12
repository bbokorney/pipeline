package main

func errorResponse(msg string) errorMessage {
	return errorMessage{
		Message: msg,
	}
}

type errorMessage struct {
	Message string `json:"message"`
}

func isValidationError(err error) bool {
	_, isType := err.(ValidationError)
	return isType
}
