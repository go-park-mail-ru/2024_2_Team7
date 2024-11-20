package auth

import (
	"context"
	"time"

	pb "kudago/internal/auth/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAPI) CreateSession(ctx context.Context, in *pb.CreateSessionRequest) (*pb.CreateSessionResponse, error) {
	session, err := s.sessionManager.CreateSession(ctx, int(in.ID))
	if err != nil {
		s.logger.Error(ctx, "create session", err)
		return nil, status.Error(codes.Internal, errInternal)
	}

	user := &pb.CreateSessionResponse{
		UserID:  int32(session.UserID),
		Token:   session.Token,
		Expires: session.Expires.Format(time.RFC3339),
	}

	return user, nil
}
