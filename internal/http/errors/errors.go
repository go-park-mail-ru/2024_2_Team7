package httpErrors

type HttpError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

var (
	ErrUserIsAuthorized = &HttpError{
		Message: "User is authorized",
		Code:    "already_authorized",
	}

	ErrUserAlreadyLoggedIn = &HttpError{
		Message: "Already logged in",
		Code:    "already_logged",
	}

	ErrUnauthorized = &HttpError{
		Message: "Unauthorized",
		Code:    "forbidden",
	}

	ErrWrongCredentials = &HttpError{
		Message: "Wrong username or password",
		Code:    "wrong_credentials",
	}

	ErrInvalidData = &HttpError{
		Message: "Can't decode JSON",
		Code:    "invalid_data",
	}

	ErrInvalidTime = &HttpError{
		Message: "Can't decode time to JSON",
		Code:    "invalid_time",
	}

	ErrInvalidCategory = &HttpError{
		Message: "Wrong or empty category",
		Code:    "invalid_category",
	}

	ErrUsernameIsAlredyTaken = &HttpError{
		Message: "Username is already taken",
		Code:    "already_taken",
	}

	ErrEmailIsAlredyTaken = &HttpError{
		Message: "Email is already taken",
		Code:    "already_taken",
	}

	ErrInternal = &HttpError{
		Message: "Internal server Error",
		Code:    "internal_Error",
	}

	ErrEventNotFound = &HttpError{
		Message: "Event not found",
		Code:    "not_found",
	}

	ErrUserNotFound = &HttpError{
		Message: "User not found",
		Code:    "not_found",
	}

	ErrAccessDenied = &HttpError{
		Message: "User doesn't own this event",
		Code:    "access_denied",
	}
)
