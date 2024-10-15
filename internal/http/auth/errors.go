package auth

type AuthError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

var (
	errUserIsAuthorized = &AuthError{
		Message: "User is authorized",
		Code:    "invalid_request",
	}

	errInvalidRequest = &AuthError{
		Message: "Invalid request",
		Code:    "invalid_request",
	}

	errInvalidData = &AuthError{
		Message: "Data is already used",
		Code:    "invalid_data",
	}

	errUserAlreadyLoggedIn = &AuthError{
		Message: "Already logged in",
		Code:    "invalid_request",
	}

	errUnauthorized = &AuthError{
		Message: "Unauthorized",
		Code:    "forbidden",
	}
	errInvalidFields = &AuthError{
		Message: "Invalid field",
		Code:    "invalid_data",
	}
)
