package events

import (
	"context"

	pb "kudago/internal/event/api"
	"kudago/internal/models"
)

func (h EventHandler) getEventsByIDs(ctx context.Context, ids []int) (map[int]models.Event, error) {
	req := &pb.GetEventsByIDsRequest{
		IDs: make([]int32, 0, len(ids)),
	}

	for _, id := range ids {
		req.IDs = append(req.IDs, int32(id))
	}

	resp, err := h.EventService.GetEventsByIDs(ctx, req)
	if err != nil {
		return nil, err
	}

	eventMap := make(map[int]models.Event, len(resp.Events))
	for _, e := range resp.Events {
		eventMap[int(e.ID)] = models.Event{
			ID:          int(e.ID),
			Title:       e.Title,
			Description: e.Description,
			EventStart:  e.EventStart,
			EventEnd:    e.EventEnd,
			ImageURL:    e.Image,
			Location:    e.Location,
			Latitude:    e.Latitude,
			Capacity:    int(e.Capacity),
			CategoryID:  int(e.CategoryID),
			AuthorID:    int(e.AuthorID),
			Tag:         e.Tag,
		}
	}

	return eventMap, nil
}
