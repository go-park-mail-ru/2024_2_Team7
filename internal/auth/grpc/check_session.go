package auth

import (
	"context"
	"errors"
	"time"

	pb "kudago/internal/auth/api"
	"kudago/internal/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAPI) CheckSession(ctx context.Context, in *pb.CheckSessionRequest) (*pb.Session, error) {
	session, err := s.sessionManager.CheckSession(ctx, in.Cookie)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, errUserNotFound)
		}
		s.logger.Error(ctx, "check session", err)
		return nil, status.Error(codes.Internal, errInternal)
	}

	user := &pb.Session{
		UserID:  int32(session.UserID),
		Token:   session.Token,
		Expires: session.Expires.Format(time.RFC3339),
	}

	return user, nil
}
