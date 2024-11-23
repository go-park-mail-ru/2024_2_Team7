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

	ErrInvalidImage = &HttpError{
		Message: "Invalid image",
		Code:    "invalid_image",
	}

	ErrInvalidImageFormat = &HttpError{
		Message: "Wrong or empty image format",
		Code:    "invalid_image",
	}

	ErrInvalidCapacity = &HttpError{
		Message: "Wrong or empty capacity",
		Code:    "invalid_capacity",
	}

	ErrInvalidID = &HttpError{
		Message: "Can't get ID",
		Code:    "invalid_id",
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

	ErrTestNotFound = &HttpError{
		Message: "Test not found",
		Code:    "not_found",
	}

	ErrAccessDenied = &HttpError{
		Message: "User doesn't own this event",
		Code:    "access_denied",
	}

	ErrBadTagLength = &HttpError{
		Message: "Tag length is limited, 20 symbols only, no empty tags",
		Code:    "invalid_tag",
	}

	ErrTooManyTags = &HttpError{
		Message: "Tags array length is limited 50 tags only",
		Code:    "invalid_tags",
	}

	ErrEventStartAfterEventEnd = &HttpError{
		Message: "Event start should be before event end",
		Code:    "invalid_time",
	}

	ErrBadEventTiming = &HttpError{
		Message: "Event start should not be in the past and event end should be before 2030",
		Code:    "invalid_time",
	}

	ErrSelfSubscription = &HttpError{
		Message: "Can't subscribe same user",
		Code:    "invalid_id",
	}

	ErrSubscriptionAlreadyExists = &HttpError{
		Message: "Already subscribed",
		Code:    "already_subscribed",
	}

	ErrSubscriptionNotFound = &HttpError{
		Message: "No subscription to delete",
		Code:    "no_subscription",
	}
)
