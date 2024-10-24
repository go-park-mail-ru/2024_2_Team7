package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	httpErrors "kudago/internal/http/errors"
	"kudago/internal/http/utils"
	"kudago/internal/models"

	"github.com/asaskevich/govalidator"
)

type AuthHandler struct {
	Service AuthService
}

type AuthService interface {
	CheckSession(ctx context.Context, cookie string) (*models.Session, bool)
	GetUserByID(ctx context.Context, ID int) (models.User, error)
	CheckCredentials(ctx context.Context, creds models.Credentials) (models.User, error)
	Register(ctx context.Context, user models.User) (models.User, error)
	CreateSession(ctx context.Context, ID int) *models.Session
	DeleteSession(ctx context.Context, token string)
}

type RegisterRequest struct {
	Username string `json:"username" valid:"required,alphanum,length(3|50)"`
	Email    string `json:"email" valid:"email,required"`
	Password string `json:"password" valid:"password,required,length(3|50)"`
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
		Service: s,
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
	err := json.NewDecoder(r.Body).Decode(&req)
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
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	user, err = h.Service.Register(r.Context(), user)
	if err != nil {
		var authErr *models.AuthError
		if errors.As(err, &authErr) {
			utils.WriteResponse(w, http.StatusConflict, authErr)
			return
		}
		fmt.Println(err)
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

	user, err := h.Service.CheckCredentials(r.Context(), creds)
	if err != nil {
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

	h.Service.DeleteSession(r.Context(), session.Token)

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

	user, err := h.Service.GetUserByID(r.Context(), session.UserID)
	fmt.Println(user, err)
	if err != nil {
		utils.WriteResponse(w, http.StatusNotFound, httpErrors.ErrUserNotFound)
		return
	}
	userResponse := userToUserResponse(user)

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

	user, err := h.Service.GetUserByID(r.Context(), session.UserID)
	if err != nil {
		utils.WriteResponse(w, http.StatusNotFound, httpErrors.ErrUserNotFound)
	}
	userResponse := userToProfileResponse(user)

	utils.WriteResponse(w, http.StatusOK, userResponse)
}

func (h *AuthHandler) setSessionCookie(w http.ResponseWriter, r *http.Request, ID int) {
	session := h.Service.CreateSession(r.Context(), ID)
	http.SetCookie(w, &http.Cookie{
		Name:     models.SessionToken,
		Value:    session.Token,
		Expires:  session.Expires,
		HttpOnly: true,
	})
}

func userToUserResponse(user models.User) UserResponse {
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
