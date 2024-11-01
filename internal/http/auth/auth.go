//go:generate mockgen -source ./auth.go -destination=./mocks/auth.go -package=mock_auth

package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	httpErrors "kudago/internal/http/errors"
	"kudago/internal/http/utils"
	"kudago/internal/models"

	"github.com/asaskevich/govalidator"
)

type AuthHandler struct {
	service AuthService
}

type AuthService interface {
	CheckSession(ctx context.Context, cookie string) (models.Session, error)
	GetUserByID(ctx context.Context, ID int) (models.User, error)
	UpdateUser(ctx context.Context, data models.NewUserData) (models.User, error)
	CheckCredentials(ctx context.Context, creds models.Credentials) (models.User, error)
	Register(ctx context.Context, registerDTO models.NewUserData) (models.User, error)
	CreateSession(ctx context.Context, ID int) (models.Session, error)
	DeleteSession(ctx context.Context, token string) error
}

type RegisterRequest struct {
	Username string `json:"username" valid:"required,alphanum,length(3|50)"`
	Email    string `json:"email" valid:"email,required"`
	Password string `json:"password" valid:"password,required,length(3|50)"`
}

type UpdateRequest struct {
	Username string `json:"username" valid:"alphanum,length(3|50)"`
	Email    string `json:"email" valid:"email"`
}

type AuthRequest struct {
	Username string `json:"username" valid:"required,alphanum,length(3|50)"`
	Password string `json:"password" valid:"password,required,length(3|50)"`
}

type AuthResponse struct {
	User UserResponse `json:"user"`
}

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	ImageURL string `json:"image"`
}

type ProfileResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	ImageURL string `json:"image"`
}

var validPasswordRegex = regexp.MustCompile(`^[a-zA-Z0-9+\-*/.;=\]\[\}\{\?]+$`)

func init() {
	govalidator.TagMap["password"] = govalidator.Validator(func(str string) bool {
		return validPasswordRegex.MatchString(str)
	})
}

func NewAuthHandler(s AuthService) *AuthHandler {
	return &AuthHandler{
		service: s,
	}
}

// @Summary Регистрация пользователя
// @Description Создает нового пользователя
// @Tags auth
// @Accept  json
// @Produce  json
// @Success 201 {object} UserResponse
// @Failure 400 {object} httpErrors.HttpError "Invalid Data / Username or Email already taken"
// @Failure 401 {object} utils.ValidationErrResponse "Validation error"
// @Failure 500 {object} httpErrors.HttpError "Internal Server Error"
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
		utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidData)
		return
	}

	_, err = govalidator.ValidateStruct(&req)
	if err != nil {
		utils.ProcessValidationErrors(w, err)
		return
	}

	// TODO errors
	r.ParseMultipartForm(1 << 20)
	file, header, err := r.FormFile("image")
	if err != nil {
		if err != http.ErrMissingFile {
			utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidData)
			return
		}
	} else {
		err = utils.GenerateFilename(header)
		if err != nil {
			utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidImage)
			return
		}
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	registerDTO := models.NewUserData{
		User:   user,
		Header: header,
		File:   &file,
	}

	user, err = h.service.Register(r.Context(), registerDTO)
	if err != nil {
		if errors.Is(err, &models.AuthError{}) {
			utils.WriteResponse(w, http.StatusConflict, err)
			return
		}
		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		return
	}

	h.setSessionCookie(w, r, user.ID)

	userResponse := UserToUserResponse(user)

	resp := AuthResponse{
		User: userResponse,
	}

	utils.WriteResponse(w, http.StatusCreated, resp)
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
		utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidData)
		return
	}

	_, err := govalidator.ValidateStruct(&req)
	if err != nil {
		utils.ProcessValidationErrors(w, err)
		return
	}

	creds := models.Credentials{
		Username: req.Username,
		Password: req.Password,
	}

	user, err := h.service.CheckCredentials(r.Context(), creds)
	if err != nil {
		utils.WriteResponse(w, http.StatusForbidden, httpErrors.ErrWrongCredentials)
		return
	}

	h.setSessionCookie(w, r, user.ID)
	userResponse := UserToUserResponse(user)

	resp := AuthResponse{
		User: userResponse,
	}

	utils.WriteResponse(w, http.StatusOK, resp)
	return
}

// @Summary Выход из системы
// @Description Выход из аккаунта
// @Tags auth
// @Success 200
// @Failure 403 {object} httpErrors.HttpError "Forbidden"
// @Router /logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	session, ok := utils.GetSessionFromContext(r.Context())
	if !ok {
		utils.WriteResponse(w, http.StatusForbidden, httpErrors.ErrUnauthorized)
		return
	}

	err := h.service.DeleteSession(r.Context(), session.Token)
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   models.SessionToken,
		MaxAge: -1, // Устанавливаем истекшее время, чтобы удалить cookie
	})

	w.WriteHeader(http.StatusOK)
}

// @Summary Проверка сессии
// @Description Возвращает информацию о пользователе, если сессия активна
// @Tags auth
// @Success 200 {object} AuthResponse
// @Router /session [get]
func (h *AuthHandler) CheckSession(w http.ResponseWriter, r *http.Request) {
	session, ok := utils.GetSessionFromContext(r.Context())
	if !ok {
		utils.WriteResponse(w, http.StatusOK, httpErrors.ErrUnauthorized)
		return
	}

	user, err := h.service.GetUserByID(r.Context(), session.UserID)
	if err != nil {
		utils.WriteResponse(w, http.StatusNotFound, httpErrors.ErrUserNotFound)
		return
	}
	userResponse := UserToUserResponse(user)

	resp := AuthResponse{
		User: userResponse,
	}

	utils.WriteResponse(w, http.StatusOK, resp)
}

// @Summary Профиль пользователя
// @Description Возвращает информацию о профиле текущего пользователя
// @Tags profile
// @Success 200 {object} ProfileResponse
// @Failure 401 {object} httpErrors.HttpError "Unauthorized"
// @Router /profile [get]
func (h *AuthHandler) Profile(w http.ResponseWriter, r *http.Request) {
	session, ok := utils.GetSessionFromContext(r.Context())
	if !ok {
		utils.WriteResponse(w, http.StatusUnauthorized, httpErrors.ErrUnauthorized)
		return
	}
	user, err := h.service.GetUserByID(r.Context(), session.UserID)
	if err != nil {
		utils.WriteResponse(w, http.StatusNotFound, httpErrors.ErrUserNotFound)
		return
	}

	userResponse := userToProfileResponse(user)

	utils.WriteResponse(w, http.StatusOK, userResponse)
}

// @Summary Профиль пользователя
// @Description Обновление информации о профиле текущего пользователя(обновление аватарки)
// @Tags profile
// @Success 200 {object} ProfileResponse
// @Failure 401 {object} httpErrors.HttpError "Unauthorized"
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

	// TODO errors
	r.ParseMultipartForm(1 << 20)
	file, header, err := r.FormFile("image")
	if err != nil {
		if err != http.ErrMissingFile {
			utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidData)
			return
		}
	} else {
		err = utils.GenerateFilename(header)
		if err != nil {
			utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidImage)
			return
		}
	}

	user := models.User{
		ID:       session.UserID,
		Username: req.Username,
		Email:    req.Email,
	}

	data := models.NewUserData{
		User:   user,
		Header: header,
		File:   &file,
	}

	newUserData, err := h.service.UpdateUser(r.Context(), data)
	if err != nil {
		utils.WriteResponse(w, http.StatusNotFound, httpErrors.ErrUserNotFound)
		return
	}
	userResponse := userToProfileResponse(newUserData)

	utils.WriteResponse(w, http.StatusOK, userResponse)
}

func (h *AuthHandler) setSessionCookie(w http.ResponseWriter, r *http.Request, ID int) {
	session, err := h.service.CreateSession(r.Context(), ID)
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     models.SessionToken,
		Value:    session.Token,
		Expires:  session.Expires,
		HttpOnly: true,
	})
}

func UserToUserResponse(user models.User) UserResponse {
	return UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		ImageURL: user.ImageURL,
	}
}

func userToProfileResponse(user models.User) ProfileResponse {
	return ProfileResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		ImageURL: user.ImageURL,
	}
}

func (h *AuthHandler) CheckSessionMiddleware(ctx context.Context, cookie string) (models.Session, error) {
	return h.service.CheckSession(ctx, cookie)
}
