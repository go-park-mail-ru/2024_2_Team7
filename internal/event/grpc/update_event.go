package grpc

import (
	"context"
	"errors"

	pb "kudago/internal/event/api"
	"kudago/internal/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAPI) UpdateEvent(ctx context.Context, req *pb.Event) (*pb.Event, error) {
	newEvent := eventPBToEvent(req)

	eventData, err := s.service.UpdateEvent(ctx, newEvent)
	if err != nil {
		s.logger.Error(ctx, "update event", err)
		switch {
		case errors.Is(err, models.ErrEventNotFound):
			return nil, status.Error(codes.NotFound, ErrEventNotFound)
		case errors.Is(err, models.ErrAccessDenied):
			return nil, status.Error(codes.PermissionDenied, ErrPermissionDenied)
		default:
			return nil, status.Error(codes.Internal, ErrInternal)
		}
	}

	event := eventToEventPB(eventData)

	return event, nil
}
