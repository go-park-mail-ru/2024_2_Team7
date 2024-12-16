package grpc

import (
	"context"
	"errors"

	pb "kudago/internal/event/api"
	"kudago/internal/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAPI) GetEventsByIDs(ctx context.Context, req *pb.GetEventsByIDsRequest) (*pb.Events, error) {
	ids := make([]int, 0, len(req.IDs))
	for _, id := range req.IDs {
		ids = append(ids, int(id))
	}

	events, err := s.getter.GetEventsByIDs(ctx, ids)
	if err != nil {
		if errors.Is(err, models.ErrEventNotFound) {
			return nil, status.Error(codes.NotFound, ErrEventNotFound)
		}
		s.logger.Error(ctx, "get events by ids", err)
		return nil, status.Error(codes.Internal, ErrInternal)
	}

	resp := writeEventsResponse(events, len(req.IDs))

	return resp, nil
}
