package handlers

import (
	"net/http"

	pb "kudago/internal/auth/api"
	httpErrors "kudago/internal/gateway/errors"
	"kudago/internal/gateway/utils"

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

	user, err := h.AuthService.GetUser(r.Context(), getUserRequest)
	if err != nil {
		h.logger.Error(r.Context(), "check session", err)
		st, ok := grpcStatus.FromError(err)
		if ok {
			switch st.Code() {
			case grpcCodes.NotFound:
				utils.WriteResponse(w, http.StatusForbidden, httpErrors.ErrUserNotFound)
				return
			case grpcCodes.Internal:
				h.logger.Error(r.Context(), "check session", st.Err())
				utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
				return
			default:
				h.logger.Error(r.Context(), "check session", st.Err())
				utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidData)
				return
			}
		}
		h.logger.Error(r.Context(), "check session", err)
		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		return
	}

	resp := userToUserResponse(user)

	utils.WriteResponse(w, http.StatusOK, resp)
	return
}
