package grpc

import (
	"context"

	pb "kudago/internal/event/api"
	"kudago/internal/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAPI) GetCategories(ctx context.Context, req *pb.Empty) (*pb.GetCategoriesResponse, error) {
	categories, err := s.getter.GetCategories(ctx)
	if err != nil {
		s.logger.Error(ctx, "get categories", err)
		return nil, status.Error(codes.Internal, errInternal)
	}
	resp := writeCategoriesResponse(categories)
	return resp, nil
}

func writeCategoriesResponse(categories []models.Category) *pb.GetCategoriesResponse {
	pbCategories := make([]*pb.Category, 0, 16)
	for _, category := range categories {
		pbCategories = append(pbCategories, &pb.Category{
			ID:   int32(category.ID),
			Name: category.Name,
		})
	}

	return &pb.GetCategoriesResponse{
		Categories: pbCategories,
	}
}
