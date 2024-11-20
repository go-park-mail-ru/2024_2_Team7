package gateway

import (
	"encoding/json"
	"fmt"
	"net/http"

	pb "kudago/internal/auth/api"
	httpErrors "kudago/internal/http/errors"

	grpcStatus "google.golang.org/grpc/status"

	"kudago/internal/http/utils"

	grpcCodes "google.golang.org/grpc/codes"

	"github.com/asaskevich/govalidator"
)

type RegisterRequest struct {
	Username string `json:"username" valid:"required,alphanum,length(3|50)"`
	Email    string `json:"email" valid:"email,required"`
	Password string `json:"password" valid:"password,required,length(3|50)"`
}

func (g *Gateway) Register(w http.ResponseWriter, r *http.Request) {
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

	registerRequest := &pb.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}

	user, err := g.authClient.Register(r.Context(), registerRequest)
	if err != nil {
		fmt.Println(err)
		st, ok := grpcStatus.FromError(err)
		if ok {
			switch st.Code() {
			case grpcCodes.AlreadyExists:
				utils.WriteResponse(w, http.StatusConflict, httpErrors.ErrUsernameIsAlredyTaken)
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
