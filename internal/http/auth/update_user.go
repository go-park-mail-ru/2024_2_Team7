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

type UpdateRequest struct {
	Username string `json:"username" valid:"alphanum,length(3|50)"`
	Email    string `json:"email" valid:"email"`
}

// UpdateUser обновляет информацию о профиле текущего пользователя, включая аватарку.
// @Summary Обновление информации о профиле пользователя
// @Description Позволяет пользователю обновить свой профиль, включая аватарку. Для этого необходимо передать JSON-объект с полями `username` и `email`, а также загрузить новый файл изображения.
// @Tags profile
// @Accept  json
// @Produce  json
// @Param json body UpdateRequest true "Данные для обновления профиля"
// @Param image formData file false "Аватарка пользователя"
// @Success 200 {object} ProfileResponse "Успешное обновление профиля"
// @Failure 400 {object} httpErrors.HttpError "Неверные данные"
// @Failure 401 {object} httpErrors.HttpError "Не авторизован"
// @Failure 404 {object} httpErrors.HttpError "Пользователь не найден"
// @Failure 409 {object} httpErrors.HttpError "Конфликт данных, пользователь с таким email или именем уже существует"
// @Router /profile [put]
func (h *AuthHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	session, ok := utils.GetSessionFromContext(r.Context())
	if !ok {
		utils.WriteResponse(w, http.StatusUnauthorized, httpErrors.ErrUnauthorized)
		return
	}

	var req UpdateRequest
	jsonData := r.FormValue("json")
	err := json.Unmarshal([]byte(jsonData), &req)
	if err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidData)
		return
	}

	_, err = govalidator.ValidateStruct(&req)
	if err != nil {
		utils.ProcessValidationErrors(w, err)
		return
	}

	user := models.User{
		ID:       session.UserID,
		Username: req.Username,
		Email:    req.Email,
	}

	media, err := utils.HandleImageUpload(r)
	if err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, err)
		return
	}

	data := models.NewUserData{
		User:  user,
		Media: media,
	}

	newUserData, err := h.service.UpdateUser(r.Context(), data)
	if err != nil {
		var authErr *models.AuthError
		if errors.As(err, &authErr) {
			utils.WriteResponse(w, http.StatusConflict, err)
			return
		}
		h.logger.Error(r.Context(), "updateUser", err)
		utils.WriteResponse(w, http.StatusNotFound, httpErrors.ErrUserNotFound)
		return
	}
	userResponse := userToUserResponse(newUserData)

	utils.WriteResponse(w, http.StatusOK, userResponse)
}
