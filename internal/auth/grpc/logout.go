package auth

import (
	"context"

	pb "kudago/internal/auth/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAPI) Logout(ctx context.Context, in *pb.LogoutRequest) (*pb.Empty, error) {
	err := s.sessionManager.DeleteSession(ctx, in.Token)
	if err != nil {
		s.logger.Error(ctx, "logout", err)
		return nil, status.Error(codes.Internal, errInternal)
	}

	return nil, nil
}
