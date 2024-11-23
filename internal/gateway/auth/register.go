package handlers

import (
	"encoding/json"
	"net/http"

	pb "kudago/internal/auth/api"
	"kudago/internal/gateway/utils"
	httpErrors "kudago/internal/http/errors"
	pbImage "kudago/internal/image/api"
	"kudago/internal/models"

	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"

	"github.com/asaskevich/govalidator"
)

type RegisterRequest struct {
	Username string `json:"username" valid:"required,alphanum,length(3|50)"`
	Email    string `json:"email" valid:"email,required"`
	Password string `json:"password" valid:"password,required,length(3|50)"`
}

func (h *AuthHandlers) Register(w http.ResponseWriter, r *http.Request) {
	_, ok := utils.GetSessionFromContext(r.Context())
	if ok {
		utils.WriteResponse(w, http.StatusForbidden, httpErrors.ErrUserIsAuthorized)
		return
	}

	req, media, reqErr := parseRegisterData(r)
	if reqErr != nil {
		utils.WriteResponse(w, http.StatusBadRequest, reqErr)
		return
	}

	_, err := govalidator.ValidateStruct(&req)
	if err != nil {
		utils.ProcessValidationErrors(w, err)
		return
	}

	url, err := h.uploadImage(r.Context(), media, w)
	if err != nil {
		return
	}

	registerRequest := &pb.RegisterRequest{
		Username:  req.Username,
		Password:  req.Password,
		Email:     req.Email,
		AvatarUrl: url,
	}

	user, err := h.AuthService.Register(r.Context(), registerRequest)
	if err != nil {
		st, ok := grpcStatus.FromError(err)
		if ok {
			switch st.Code() {
			case grpcCodes.AlreadyExists:
				utils.WriteResponse(w, http.StatusConflict, httpErrors.ErrUsernameIsAlredyTaken)
				return
			case grpcCodes.Internal:
				h.logger.Error(r.Context(), "register", st.Err())
				utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
				return
			default:
				h.logger.Error(r.Context(), "register", st.Err())
				utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidData)
				return
			}
		}

		h.logger.Error(r.Context(), "register", err)

		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		return
	}
	h.logger.Error(r.Context(), "set cookie", err)

	err = h.setSessionCookie(w, r, int(user.ID))
	if err != nil {
		h.logger.Error(r.Context(), "set cookie", err)
		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		return
	}

	resp := userToUserResponse(user)

	utils.WriteResponse(w, http.StatusOK, resp)
	return
}

func parseRegisterData(r *http.Request) (models.User, *pbImage.UploadRequest, *httpErrors.HttpError) {
	var req models.User
	jsonData := r.FormValue("json")
	err := json.Unmarshal([]byte(jsonData), &req)
	if err != nil {
		return req, nil, httpErrors.ErrInvalidData
	}

	media, err := utils.HandleImageUpload(r)
	if err != nil {
		return req, nil, httpErrors.ErrInvalidImage
	}

	return req, media, nil
}
