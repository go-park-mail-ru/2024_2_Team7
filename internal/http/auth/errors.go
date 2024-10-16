package auth

type AuthError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

var (
	errUserIsAuthorized = &AuthError{
		Message: "User is authorized",
		Code:    "already_authorized",
	}

	errUserAlreadyLoggedIn = &AuthError{
		Message: "Already logged in",
		Code:    "already_logged",
	}

	errUnauthorized = &AuthError{
		Message: "Unauthorized",
		Code:    "forbidden",
	}
	errWrongCredentials = &AuthError{
		Message: "Wrong username or password",
		Code:    "wrong_credentials",
	}

	errInvalidData = &AuthError{
		Message: "Can't decode JSON",
		Code:    "invalid_data",
	}

	errUsernameIsAlredyTaken = &AuthError{
		Message: "Username is already taken",
		Code:    "already_taken",
	}

	errEmailIsAlredyTaken = &AuthError{
		Message: "Email is already taken",
		Code:    "already_taken",
	}

	errInternal = &AuthError{
		Message: "Internal server error",
		Code:    "internal_error",
	}
)
