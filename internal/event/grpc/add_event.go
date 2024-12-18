package grpc

import (
	"context"

	pb "kudago/internal/event/api"
	"kudago/internal/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAPI) AddEvent(ctx context.Context, req *pb.Event) (*pb.Event, error) {
	newEvent := eventPBToEvent(req)

	eventData, err := s.service.AddEvent(ctx, newEvent)
	if err != nil {
		switch err {
		case models.ErrInvalidCategory:
			return nil, status.Error(codes.InvalidArgument, ErrBadData)
		}
		s.logger.Error(ctx, "add event", err)
		return nil, status.Error(codes.Internal, ErrInternal)
	}

	event := eventToEventPB(eventData)

	return event, nil
}
