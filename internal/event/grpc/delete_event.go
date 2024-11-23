package http

import (
	"context"

	pb "kudago/internal/event/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAPI) DeleteEvent(ctx context.Context, req *pb.DeleteEventRequest) (*pb.Empty, error) {
	err := s.service.DeleteEvent(ctx, int(req.EventID), int(req.AuthorID))
	if err != nil {
		s.logger.Error(ctx, "delete event", err)
		return nil, status.Error(codes.Internal, errInternal)
	}

	return nil, nil
}
