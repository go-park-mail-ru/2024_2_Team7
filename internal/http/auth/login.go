package auth

import (
	"encoding/json"
	"net/http"

	httpErrors "kudago/internal/http/errors"
	"kudago/internal/http/utils"
	"kudago/internal/models"

	"github.com/asaskevich/govalidator"
)

type AuthRequest struct {
	Username string `json:"username" valid:"required,alphanum,length(3|50)"`
	Password string `json:"password" valid:"password,required,length(3|50)"`
}

// @Summary Авторизация пользователя
// @Description Авторизует пользователя
// @Tags auth
// @Accept  json
// @Produce  json
// @Success 200 {object} UserResponse
// @Failure 400 {object} httpErrors.HttpError "Wrong Credentials"
// @Failure 401 {object} utils.ValidationErrResponse "Validation error"
// @Failure 403 {object} httpErrors.HttpError "User is alredy logged in"
// @Failure 500 {object} httpErrors.HttpError "Internal Server Error"
// @Router /login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	_, ok := utils.GetSessionFromContext(r.Context())
	if ok {
		utils.WriteResponse(w, http.StatusForbidden, httpErrors.ErrUserAlreadyLoggedIn)
		return
	}

	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(r.Context(), "login", err)
		utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidData)
		return
	}

	_, err := govalidator.ValidateStruct(&req)
	if err != nil {
		h.logger.Error(r.Context(), "login", err)
		utils.ProcessValidationErrors(w, err)
		return
	}

	creds := models.Credentials{
		Username: req.Username,
		Password: req.Password,
	}

	user, err := h.service.CheckCredentials(r.Context(), creds)
	if err != nil {
		h.logger.Error(r.Context(), "login", err)
		utils.WriteResponse(w, http.StatusForbidden, httpErrors.ErrWrongCredentials)
		return
	}

	h.setSessionCookie(w, r, user.ID)
	userResponse := userToUserResponse(user)

	resp := AuthResponse{
		User: userResponse,
	}

	utils.WriteResponse(w, http.StatusOK, resp)
	return
}
