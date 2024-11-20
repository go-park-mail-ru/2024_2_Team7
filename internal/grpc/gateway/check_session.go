package gateway

import (
	"net/http"

	pb "kudago/internal/auth/api"
	httpErrors "kudago/internal/http/errors"
	"kudago/internal/http/utils"
)

func (g *Gateway) CheckSession(w http.ResponseWriter, r *http.Request) {
	session, ok := utils.GetSessionFromContext(r.Context())
	if !ok {
		utils.WriteResponse(w, http.StatusOK, httpErrors.ErrUnauthorized)
		return
	}

	sessionRequest := &pb.CheckSessionRequest{
		ID: int32(session.UserID),
	}

	user, err := g.authClient.CheckSession(r.Context(), sessionRequest)
	if err != nil {
		utils.WriteResponse(w, http.StatusNotFound, httpErrors.ErrUserNotFound)
		return
	}

	resp := userToUserResponse(user)

	utils.WriteResponse(w, http.StatusOK, resp)
	return
}
