package grpc

import (
	"context"

	pb "kudago/internal/event/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAPI) GetFavorites(ctx context.Context, req *pb.GetFavoritesRequest) (*pb.Events, error) {
	params := getPaginationParams(req.Params)
	eventsData, err := s.getter.GetFavorites(ctx, int(req.UserID), params)
	if err != nil {
		s.logger.Error(ctx, "get favorites", err)
		return nil, status.Error(codes.Internal, ErrInternal)
	}

	event := writeEventsResponse(eventsData, params.Limit)

	return event, nil
}
