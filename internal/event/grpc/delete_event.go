package grpc

import (
	"context"
	"errors"

	pb "kudago/internal/event/api"
	"kudago/internal/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAPI) DeleteEvent(ctx context.Context, req *pb.DeleteEventRequest) (*pb.Empty, error) {
	err := s.service.DeleteEvent(ctx, int(req.EventID), int(req.AuthorID))
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, status.Error(codes.NotFound, ErrEventNotFound)
		}
		s.logger.Error(ctx, "delete event", err)
		return nil, status.Error(codes.Internal, ErrInternal)
	}
	return nil, nil
}
