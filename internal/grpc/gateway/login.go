package gateway

import (
	"encoding/json"
	"net/http"

	pb "kudago/internal/auth/api"
	httpErrors "kudago/internal/http/errors"
	"kudago/internal/http/utils"

	"github.com/asaskevich/govalidator"
	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

type LoginRequest struct {
	Username string `json:"username" valid:"required,alphanum,length(3|50)"`
	Password string `json:"password" valid:"password,required,length(3|50)"`
}

func (g *Gateway) Login(w http.ResponseWriter, r *http.Request) {
	_, ok := utils.GetSessionFromContext(r.Context())
	if ok {
		utils.WriteResponse(w, http.StatusForbidden, httpErrors.ErrUserAlreadyLoggedIn)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidData)
		return
	}

	_, err := govalidator.ValidateStruct(&req)
	if err != nil {
		utils.ProcessValidationErrors(w, err)
		return
	}

	creds := &pb.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}

	user, err := g.authClient.Login(r.Context(), creds)
	if err != nil {
		st, ok := grpcStatus.FromError(err)
		if ok {
			switch st.Code() {
			case grpcCodes.NotFound:
				utils.WriteResponse(w, http.StatusForbidden, httpErrors.ErrWrongCredentials)
				return
			case grpcCodes.Internal:
				utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
				return
			default:
				utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidData)
				return
			}
		}

		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		return
	}

	err = g.setSessionCookie(w, r, int(user.ID))
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		return
	}
	resp := userToUserResponse(user)

	utils.WriteResponse(w, http.StatusOK, resp)
	return
}
