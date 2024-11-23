package http

import (
	"context"

	pb "kudago/internal/event/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAPI) AddEventToFavorites(ctx context.Context, req *pb.FavoriteEvent) (*pb.Empty, error) {
	newFavorite := favoritePBToFavorite(req)

	err := s.service.AddEventToFavorites(ctx, newFavorite)
	if err != nil {
		s.logger.Error(ctx, "add event to favorites", err)
		return nil, status.Error(codes.Internal, errInternal)
	}

	return nil, nil
}
