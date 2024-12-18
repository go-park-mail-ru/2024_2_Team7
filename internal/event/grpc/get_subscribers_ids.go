package grpc

import (
	"context"

	pb "kudago/internal/event/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAPI) GetSubscribersIDs(ctx context.Context, req *pb.GetSubscribersIDsRequest) (*pb.GetUserIDsResponse, error) {
	ids, err := s.getter.GetSubscribersIDs(ctx, int(req.UserID))
	if err != nil {
		s.logger.Error(ctx, "get subscribers", err)
		return nil, status.Error(codes.Internal, ErrInternal)
	}

	resp := &pb.GetUserIDsResponse{
		IDs: make([]int32, 0, len(ids)),
	}

	for _, id := range ids {
		resp.IDs = append(resp.IDs, int32(id))
	}

	return resp, nil
}
