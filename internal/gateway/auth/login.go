package handlers

import (
	"io"
	"net/http"

	pb "kudago/internal/auth/api"
	httpErrors "kudago/internal/gateway/errors"
	"kudago/internal/gateway/utils"

	"github.com/asaskevich/govalidator"
	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	_, ok := utils.GetSessionFromContext(r.Context())
	if ok {
		utils.WriteResponse(w, http.StatusForbidden, httpErrors.ErrUserAlreadyLoggedIn)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidData)
		return
	}
	defer r.Body.Close()

	var req LoginRequest
	err = req.UnmarshalJSON([]byte(body))
	if err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidData)
		return
	}

	_, err = govalidator.ValidateStruct(&req)
	if err != nil {
		utils.ProcessValidationErrors(w, err)
		return
	}

	creds := &pb.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}

	user, err := h.AuthService.Login(r.Context(), creds)
	if err != nil {
		st, ok := grpcStatus.FromError(err)
		if ok {
			switch st.Code() {
			case grpcCodes.NotFound:
				utils.WriteResponse(w, http.StatusForbidden, httpErrors.ErrWrongCredentials)
				return
			}
		}
		h.logger.Error(r.Context(), "login", err)
		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		return
	}

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
