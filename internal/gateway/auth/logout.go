package handlers

import (
	"net/http"

	pb "kudago/internal/auth/api"
	httpErrors "kudago/internal/gateway/errors"
	"kudago/internal/gateway/utils"
	"kudago/internal/models"
)

func (h *AuthHandlers) Logout(w http.ResponseWriter, r *http.Request) {
	session, ok := utils.GetSessionFromContext(r.Context())
	if !ok {
		utils.WriteResponse(w, http.StatusForbidden, httpErrors.ErrUnauthorized)
		return
	}

	req := &pb.LogoutRequest{
		Token: session.Token,
	}

	_, err := h.AuthService.Logout(r.Context(), req)
	if err != nil {
		h.logger.Error(r.Context(), "logout", err)

		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   models.SessionToken,
		MaxAge: -1, // Устанавливаем истекшее время, чтобы удалить cookie
	})

	w.WriteHeader(http.StatusOK)
}
