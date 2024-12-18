package grpc

import (
	"context"

	pb "kudago/internal/image/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAPI) DeleteImage(ctx context.Context, req *pb.DeleteRequest) (*pb.Empty, error) {
	err := s.service.DeleteImage(ctx, req.FileUrl)
	if err != nil {
		s.logger.Error(ctx, "delete image", err)
		return nil, status.Error(codes.Internal, ErrInternal)
	}

	return nil, nil
}
