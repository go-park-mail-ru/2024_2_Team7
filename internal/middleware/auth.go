package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	pb "kudago/internal/auth/api"
	"kudago/internal/gateway/utils"
	httpErrors "kudago/internal/http/errors"
	"kudago/internal/models"

	"google.golang.org/grpc"
)

const (
	SessionToken = "session_token"
	SessionKey   = "session"
)

var whitelist = []string{
	"/login",
	"/register",
	"/events",
	"/static",
	"/session",
	"/logout",
	"/docs",
	"/categories",
	"/swagger",
	"/profile",
}

type sessionChecker interface {
	CheckSession(ctx context.Context, req *pb.CheckSessionRequest, opts ...grpc.CallOption) (*pb.Session, error)
}

func AuthMiddleware(sessionChecker sessionChecker, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(SessionToken)
		if err == nil {
			req := &pb.CheckSessionRequest{
				Cookie: cookie.Value,
			}

			sessionPB, err := sessionChecker.CheckSession(r.Context(), req)
			if err == nil {
				session := sessionPBToSession(sessionPB)
				ctx := utils.SetSessionInContext(r.Context(), session)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}
		}

		for _, path := range whitelist {
			if strings.HasPrefix(r.URL.Path, path) {
				next.ServeHTTP(w, r)
				return
			}
		}

		utils.WriteResponse(w, http.StatusUnauthorized, httpErrors.ErrUnauthorized)
		return
	})
}

func sessionPBToSession(sessionPB *pb.Session) models.Session {
	expires, _ := time.Parse(time.RFC3339, sessionPB.Expires)
	return models.Session{
		UserID:  int(sessionPB.UserID),
		Token:   sessionPB.Token,
		Expires: expires,
	}
}
