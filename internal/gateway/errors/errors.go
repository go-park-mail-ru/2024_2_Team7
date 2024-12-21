//go:generate easyjson errors.go
package httpErrors

//easyjson:json
type HttpError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

var (
	ErrUserIsAuthorized = &HttpError{
		Message: "Пользователь уже авторизован",
		Code:    "already_authorized",
	}

	ErrUserAlreadyDidTest = &HttpError{
		Message: "Пользователь уже прошел тест",
		Code:    "already_answered",
	}

	ErrEventAlreadyAddedToFavorites = &HttpError{
		Message: "Событие уже добавлено в избранное",
		Code:    "already_added",
	}

	ErrUserAlreadyLoggedIn = &HttpError{
		Message: "Пользователь уже вошел в систему",
		Code:    "already_logged",
	}

	ErrUnauthorized = &HttpError{
		Message: "Неавторизован",
		Code:    "forbidden",
	}

	ErrWrongCredentials = &HttpError{
		Message: "Неправильное имя пользователя или пароль",
		Code:    "wrong_credentials",
	}

	ErrInvalidData = &HttpError{
		Message: "Невозможно декодировать JSON",
		Code:    "invalid_data",
	}

	ErrInvalidTime = &HttpError{
		Message: "Невозможно декодировать время в JSON",
		Code:    "invalid_time",
	}

	ErrInvalidCategory = &HttpError{
		Message: "Неверная или пустая категория",
		Code:    "invalid_category",
	}

	ErrInvalidImage = &HttpError{
		Message: "Неверное изображение",
		Code:    "invalid_image",
	}

	ErrInvalidImageFormat = &HttpError{
		Message: "Неверный или пустой формат изображения",
		Code:    "invalid_image",
	}

	ErrInvalidCapacity = &HttpError{
		Message: "Неверная или пустая вместимость",
		Code:    "invalid_capacity",
	}

	ErrInvalidID = &HttpError{
		Message: "Невозможно получить ID",
		Code:    "invalid_id",
	}

	ErrUsernameIsAlredyTaken = &HttpError{
		Message: "Имя пользователя уже занято",
		Code:    "already_taken",
	}

	ErrEmailIsAlredyTaken = &HttpError{
		Message: "Электронная почта уже занята",
		Code:    "already_taken",
	}

	ErrInternal = &HttpError{
		Message: "Внутренняя ошибка сервера",
		Code:    "internal_error",
	}

	ErrEventNotFound = &HttpError{
		Message: "Событие не найдено",
		Code:    "not_found",
	}

	ErrUserNotFound = &HttpError{
		Message: "Пользователь не найден",
		Code:    "not_found",
	}

	ErrTestNotFound = &HttpError{
		Message: "Тест не найден",
		Code:    "not_found",
	}

	ErrAccessDenied = &HttpError{
		Message: "Пользователь не владеет этим событием",
		Code:    "access_denied",
	}

	ErrBadTagLength = &HttpError{
		Message: "Длина тега ограничена: максимум 20 символов, пустые теги недопустимы",
		Code:    "invalid_tag",
	}

	ErrTooManyTags = &HttpError{
		Message: "Массив тегов ограничен: максимум 50 тегов",
		Code:    "invalid_tags",
	}

	ErrEventStartAfterEventEnd = &HttpError{
		Message: "Начало события должно быть до его окончания",
		Code:    "invalid_time",
	}

	ErrBadEventTiming = &HttpError{
		Message: "Начало события не должно быть в прошлом, а окончание должно быть до 2030 года",
		Code:    "invalid_time",
	}

	ErrSelfSubscription = &HttpError{
		Message: "Нельзя подписаться на самого себя",
		Code:    "invalid_id",
	}

	ErrSubscriptionAlreadyExists = &HttpError{
		Message: "Уже подписан",
		Code:    "already_subscribed",
	}

	ErrSubscriptionNotFound = &HttpError{
		Message: "Подписка для удаления не найдена",
		Code:    "no_subscription",
	}
)
