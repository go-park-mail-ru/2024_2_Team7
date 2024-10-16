package events

type EventError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

var (
	errInternal = &EventError{
		Message: "Internal server error",
		Code:    "internal_error",
	}

	errInvalidData = &EventError{
		Message: "Can't decode JSON",
		Code:    "invalid_data",
	}

	errUnauthorized = &EventError{
		Message: "Unauthorized",
		Code:    "forbidden",
	}
	
	errEventNotFound = &EventError{
		Message: "Event not found",
		Code:    "not_found",
	}

	errAccessDenied = &EventError{
		Message: "User doesn't own this event",
		Code:    "access_denied",
	}
)
