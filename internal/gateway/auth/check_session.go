package handlers

import (
	"net/http"

	pb "kudago/internal/auth/api"
	httpErrors "kudago/internal/http/errors"
	"kudago/internal/http/utils"

	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

func (h *AuthHandlers) CheckSession(w http.ResponseWriter, r *http.Request) {
	session, ok := utils.GetSessionFromContext(r.Context())
	if !ok {
		utils.WriteResponse(w, http.StatusOK, httpErrors.ErrUnauthorized)
		return
	}

	getUserRequest := &pb.GetUserRequest{
		ID: int32(session.UserID),
	}

	user, err := h.Gateway.AuthService.GetUser(r.Context(), getUserRequest)
	if err != nil {
		st, ok := grpcStatus.FromError(err)
		if ok {
			switch st.Code() {
			case grpcCodes.NotFound:
				utils.WriteResponse(w, http.StatusForbidden, httpErrors.ErrUserNotFound)
				return
			case grpcCodes.Internal:
				h.Gateway.Logger.Error(r.Context(), "check session", st.Err())
				utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
				return
			default:
				h.Gateway.Logger.Error(r.Context(), "check session", st.Err())
				utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidData)
				return
			}
		}
		h.Gateway.Logger.Error(r.Context(), "check session", err)
		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		return
	}

	resp := userToUserResponse(user)

	utils.WriteResponse(w, http.StatusOK, resp)
	return
}
