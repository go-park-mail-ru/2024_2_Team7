package auth

import (
	"context"
	"time"

	pb "kudago/internal/auth/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAPI) CheckSession(ctx context.Context, in *pb.CheckSessionRequest) (*pb.Session, error) {
	session, err := s.sessionManager.CheckSession(ctx, in.Cookie)
	if err != nil {
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
