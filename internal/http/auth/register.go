package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	httpErrors "kudago/internal/http/errors"
	"kudago/internal/http/utils"
	"kudago/internal/models"

	"github.com/asaskevich/govalidator"
)

type RegisterRequest struct {
	Username string `json:"username" valid:"required,alphanum,length(3|50)"`
	Email    string `json:"email" valid:"email,required"`
	Password string `json:"password" valid:"password,required,length(3|50)"`
}

// Register создает нового пользователя в системе.
// @Summary Регистрация пользователя
// @Description Создает нового пользователя. Необходимо передать JSON-объект с полями `username`, `email` и `password`. Если пользователь уже авторизован, запрос будет отклонен.
// @Tags auth
// @Accept  json
// @Produce  json
// @Param json body RegisterRequest true "Данные для регистрации пользователя"
// @Param image formData file false "Аватарка пользователя"
// @Success 201 {object} AuthResponse "Успешная регистрация пользователя"
// @Failure 400 {object} httpErrors.HttpError "Неверные данные / Имя пользователя или электронная почта уже заняты"
// @Failure 401 {object} utils.ValidationErrResponse "Ошибка валидации"
// @Failure 500 {object} httpErrors.HttpError "Внутренняя ошибка сервера"
// @Router /register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	_, ok := utils.GetSessionFromContext(r.Context())
	if ok {
		utils.WriteResponse(w, http.StatusForbidden, httpErrors.ErrUserIsAuthorized)
		return
	}

	var req RegisterRequest
	jsonData := r.FormValue("json")
	err := json.Unmarshal([]byte(jsonData), &req)
	if err != nil {
		h.logger.Error(r.Context(), "register", err)
		utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidData)
		return
	}

	err = utils.SanitizeStruct(&req)
	if err != nil {
		h.logger.Error(r.Context(), "sanitize error", err)
		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		return
	}

	_, err = govalidator.ValidateStruct(&req)
	if err != nil {
		h.logger.Error(r.Context(), "register", err)
		utils.ProcessValidationErrors(w, err)
		return
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	media, err := utils.HandleImageUpload(r)
	if err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, err)
		return
	}

	registerDTO := models.NewUserData{
		User:  user,
		Media: media,
	}

	user, err = h.service.Register(r.Context(), registerDTO)
	if err != nil {
		var authErr *models.AuthError
		if errors.As(err, &authErr) {
			utils.WriteResponse(w, http.StatusConflict, err)
			return
		}
		h.logger.Error(r.Context(), "register", err)
		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		return
	}

	h.setSessionCookie(w, r, user.ID)

	userResponse := userToUserResponse(user)

	resp := AuthResponse{
		User: userResponse,
	}

	utils.WriteResponse(w, http.StatusCreated, resp)
}
