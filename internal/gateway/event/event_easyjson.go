// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package events

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
	models "kudago/internal/models"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonF642ad3eDecodeKudagoInternalGatewayEvent(in *jlexer.Lexer, out *NotificationWithEvent) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "notification":
			(out.Notification).UnmarshalEasyJSON(in)
		case "event":
			easyjsonF642ad3eDecodeKudagoInternalModels(in, &out.Event)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonF642ad3eEncodeKudagoInternalGatewayEvent(out *jwriter.Writer, in NotificationWithEvent) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"notification\":"
		out.RawString(prefix[1:])
		(in.Notification).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"event\":"
		out.RawString(prefix)
		easyjsonF642ad3eEncodeKudagoInternalModels(out, in.Event)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v NotificationWithEvent) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonF642ad3eEncodeKudagoInternalGatewayEvent(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v NotificationWithEvent) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF642ad3eEncodeKudagoInternalGatewayEvent(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *NotificationWithEvent) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonF642ad3eDecodeKudagoInternalGatewayEvent(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *NotificationWithEvent) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF642ad3eDecodeKudagoInternalGatewayEvent(l, v)
}
func easyjsonF642ad3eDecodeKudagoInternalModels(in *jlexer.Lexer, out *models.Event) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = int(in.Int())
		case "title":
			out.Title = string(in.String())
		case "description":
			out.Description = string(in.String())
		case "event_start":
			out.EventStart = string(in.String())
		case "event_finish":
			out.EventEnd = string(in.String())
		case "location":
			out.Location = string(in.String())
		case "capacity":
			out.Capacity = int(in.Int())
		case "created_at":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.CreatedAt).UnmarshalJSON(data))
			}
		case "category_id":
			out.CategoryID = int(in.Int())
		case "author":
			out.AuthorID = int(in.Int())
		case "Latitude":
			out.Latitude = float64(in.Float64())
		case "Longitude":
			out.Longitude = float64(in.Float64())
		case "tag":
			if in.IsNull() {
				in.Skip()
				out.Tag = nil
			} else {
				in.Delim('[')
				if out.Tag == nil {
					if !in.IsDelim(']') {
						out.Tag = make([]string, 0, 4)
					} else {
						out.Tag = []string{}
					}
				} else {
					out.Tag = (out.Tag)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.Tag = append(out.Tag, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "image":
			out.ImageURL = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonF642ad3eEncodeKudagoInternalModels(out *jwriter.Writer, in models.Event) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int(int(in.ID))
	}
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"event_start\":"
		out.RawString(prefix)
		out.String(string(in.EventStart))
	}
	{
		const prefix string = ",\"event_finish\":"
		out.RawString(prefix)
		out.String(string(in.EventEnd))
	}
	{
		const prefix string = ",\"location\":"
		out.RawString(prefix)
		out.String(string(in.Location))
	}
	{
		const prefix string = ",\"capacity\":"
		out.RawString(prefix)
		out.Int(int(in.Capacity))
	}
	{
		const prefix string = ",\"created_at\":"
		out.RawString(prefix)
		out.Raw((in.CreatedAt).MarshalJSON())
	}
	{
		const prefix string = ",\"category_id\":"
		out.RawString(prefix)
		out.Int(int(in.CategoryID))
	}
	{
		const prefix string = ",\"author\":"
		out.RawString(prefix)
		out.Int(int(in.AuthorID))
	}
	{
		const prefix string = ",\"Latitude\":"
		out.RawString(prefix)
		out.Float64(float64(in.Latitude))
	}
	{
		const prefix string = ",\"Longitude\":"
		out.RawString(prefix)
		out.Float64(float64(in.Longitude))
	}
	{
		const prefix string = ",\"tag\":"
		out.RawString(prefix)
		if in.Tag == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Tag {
				if v2 > 0 {
					out.RawByte(',')
				}
				out.String(string(v3))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"image\":"
		out.RawString(prefix)
		out.String(string(in.ImageURL))
	}
	out.RawByte('}')
}
func easyjsonF642ad3eDecodeKudagoInternalGatewayEvent1(in *jlexer.Lexer, out *NewEventResponse) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "event":
			(out.Event).UnmarshalEasyJSON(in)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonF642ad3eEncodeKudagoInternalGatewayEvent1(out *jwriter.Writer, in NewEventResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"event\":"
		out.RawString(prefix[1:])
		(in.Event).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v NewEventResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonF642ad3eEncodeKudagoInternalGatewayEvent1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v NewEventResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF642ad3eEncodeKudagoInternalGatewayEvent1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *NewEventResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonF642ad3eDecodeKudagoInternalGatewayEvent1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *NewEventResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF642ad3eDecodeKudagoInternalGatewayEvent1(l, v)
}
func easyjsonF642ad3eDecodeKudagoInternalGatewayEvent2(in *jlexer.Lexer, out *NewEventRequest) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "title":
			out.Title = string(in.String())
		case "description":
			out.Description = string(in.String())
		case "location":
			out.Location = string(in.String())
		case "category_id":
			out.Category = int(in.Int())
		case "capacity":
			out.Capacity = int(in.Int())
		case "tag":
			if in.IsNull() {
				in.Skip()
				out.Tag = nil
			} else {
				in.Delim('[')
				if out.Tag == nil {
					if !in.IsDelim(']') {
						out.Tag = make([]string, 0, 4)
					} else {
						out.Tag = []string{}
					}
				} else {
					out.Tag = (out.Tag)[:0]
				}
				for !in.IsDelim(']') {
					var v4 string
					v4 = string(in.String())
					out.Tag = append(out.Tag, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "event_start":
			out.EventStart = string(in.String())
		case "event_end":
			out.EventEnd = string(in.String())
		case "Latitude":
			out.Latitude = float64(in.Float64())
		case "Longitude":
			out.Longitude = float64(in.Float64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonF642ad3eEncodeKudagoInternalGatewayEvent2(out *jwriter.Writer, in NewEventRequest) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix[1:])
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"location\":"
		out.RawString(prefix)
		out.String(string(in.Location))
	}
	{
		const prefix string = ",\"category_id\":"
		out.RawString(prefix)
		out.Int(int(in.Category))
	}
	{
		const prefix string = ",\"capacity\":"
		out.RawString(prefix)
		out.Int(int(in.Capacity))
	}
	{
		const prefix string = ",\"tag\":"
		out.RawString(prefix)
		if in.Tag == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range in.Tag {
				if v5 > 0 {
					out.RawByte(',')
				}
				out.String(string(v6))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"event_start\":"
		out.RawString(prefix)
		out.String(string(in.EventStart))
	}
	{
		const prefix string = ",\"event_end\":"
		out.RawString(prefix)
		out.String(string(in.EventEnd))
	}
	{
		const prefix string = ",\"Latitude\":"
		out.RawString(prefix)
		out.Float64(float64(in.Latitude))
	}
	{
		const prefix string = ",\"Longitude\":"
		out.RawString(prefix)
		out.Float64(float64(in.Longitude))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v NewEventRequest) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonF642ad3eEncodeKudagoInternalGatewayEvent2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v NewEventRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF642ad3eEncodeKudagoInternalGatewayEvent2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *NewEventRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonF642ad3eDecodeKudagoInternalGatewayEvent2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *NewEventRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF642ad3eDecodeKudagoInternalGatewayEvent2(l, v)
}
func easyjsonF642ad3eDecodeKudagoInternalGatewayEvent3(in *jlexer.Lexer, out *InviteNotificationRequest) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "user_id":
			out.UserID = int(in.Int())
		case "event_id":
			out.EventID = int(in.Int())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonF642ad3eEncodeKudagoInternalGatewayEvent3(out *jwriter.Writer, in InviteNotificationRequest) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"user_id\":"
		out.RawString(prefix[1:])
		out.Int(int(in.UserID))
	}
	{
		const prefix string = ",\"event_id\":"
		out.RawString(prefix)
		out.Int(int(in.EventID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v InviteNotificationRequest) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonF642ad3eEncodeKudagoInternalGatewayEvent3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v InviteNotificationRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF642ad3eEncodeKudagoInternalGatewayEvent3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *InviteNotificationRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonF642ad3eDecodeKudagoInternalGatewayEvent3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *InviteNotificationRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF642ad3eDecodeKudagoInternalGatewayEvent3(l, v)
}
func easyjsonF642ad3eDecodeKudagoInternalGatewayEvent4(in *jlexer.Lexer, out *GetNotificationsResponse) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "notifications":
			if in.IsNull() {
				in.Skip()
				out.Notifications = nil
			} else {
				in.Delim('[')
				if out.Notifications == nil {
					if !in.IsDelim(']') {
						out.Notifications = make([]NotificationWithEvent, 0, 0)
					} else {
						out.Notifications = []NotificationWithEvent{}
					}
				} else {
					out.Notifications = (out.Notifications)[:0]
				}
				for !in.IsDelim(']') {
					var v7 NotificationWithEvent
					(v7).UnmarshalEasyJSON(in)
					out.Notifications = append(out.Notifications, v7)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonF642ad3eEncodeKudagoInternalGatewayEvent4(out *jwriter.Writer, in GetNotificationsResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"notifications\":"
		out.RawString(prefix[1:])
		if in.Notifications == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v8, v9 := range in.Notifications {
				if v8 > 0 {
					out.RawByte(',')
				}
				(v9).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v GetNotificationsResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonF642ad3eEncodeKudagoInternalGatewayEvent4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v GetNotificationsResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF642ad3eEncodeKudagoInternalGatewayEvent4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *GetNotificationsResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonF642ad3eDecodeKudagoInternalGatewayEvent4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *GetNotificationsResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF642ad3eDecodeKudagoInternalGatewayEvent4(l, v)
}
func easyjsonF642ad3eDecodeKudagoInternalGatewayEvent5(in *jlexer.Lexer, out *GetEventsResponse) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "events":
			if in.IsNull() {
				in.Skip()
				out.Events = nil
			} else {
				in.Delim('[')
				if out.Events == nil {
					if !in.IsDelim(']') {
						out.Events = make([]EventResponse, 0, 0)
					} else {
						out.Events = []EventResponse{}
					}
				} else {
					out.Events = (out.Events)[:0]
				}
				for !in.IsDelim(']') {
					var v10 EventResponse
					(v10).UnmarshalEasyJSON(in)
					out.Events = append(out.Events, v10)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonF642ad3eEncodeKudagoInternalGatewayEvent5(out *jwriter.Writer, in GetEventsResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"events\":"
		out.RawString(prefix[1:])
		if in.Events == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v11, v12 := range in.Events {
				if v11 > 0 {
					out.RawByte(',')
				}
				(v12).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v GetEventsResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonF642ad3eEncodeKudagoInternalGatewayEvent5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v GetEventsResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF642ad3eEncodeKudagoInternalGatewayEvent5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *GetEventsResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonF642ad3eDecodeKudagoInternalGatewayEvent5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *GetEventsResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF642ad3eDecodeKudagoInternalGatewayEvent5(l, v)
}
func easyjsonF642ad3eDecodeKudagoInternalGatewayEvent6(in *jlexer.Lexer, out *GetCategoriesResponse) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "categories":
			if in.IsNull() {
				in.Skip()
				out.Categories = nil
			} else {
				in.Delim('[')
				if out.Categories == nil {
					if !in.IsDelim(']') {
						out.Categories = make([]models.Category, 0, 2)
					} else {
						out.Categories = []models.Category{}
					}
				} else {
					out.Categories = (out.Categories)[:0]
				}
				for !in.IsDelim(']') {
					var v13 models.Category
					easyjsonF642ad3eDecodeKudagoInternalModels1(in, &v13)
					out.Categories = append(out.Categories, v13)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonF642ad3eEncodeKudagoInternalGatewayEvent6(out *jwriter.Writer, in GetCategoriesResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"categories\":"
		out.RawString(prefix[1:])
		if in.Categories == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v14, v15 := range in.Categories {
				if v14 > 0 {
					out.RawByte(',')
				}
				easyjsonF642ad3eEncodeKudagoInternalModels1(out, v15)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v GetCategoriesResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonF642ad3eEncodeKudagoInternalGatewayEvent6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v GetCategoriesResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF642ad3eEncodeKudagoInternalGatewayEvent6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *GetCategoriesResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonF642ad3eDecodeKudagoInternalGatewayEvent6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *GetCategoriesResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF642ad3eDecodeKudagoInternalGatewayEvent6(l, v)
}
func easyjsonF642ad3eDecodeKudagoInternalModels1(in *jlexer.Lexer, out *models.Category) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = int(in.Int())
		case "name":
			out.Name = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonF642ad3eEncodeKudagoInternalModels1(out *jwriter.Writer, in models.Category) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int(int(in.ID))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	out.RawByte('}')
}
func easyjsonF642ad3eDecodeKudagoInternalGatewayEvent7(in *jlexer.Lexer, out *EventResponse) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = int(in.Int())
		case "title":
			out.Title = string(in.String())
		case "description":
			out.Description = string(in.String())
		case "location":
			out.Location = string(in.String())
		case "category_id":
			out.Category = int(in.Int())
		case "capacity":
			out.Capacity = int(in.Int())
		case "tag":
			if in.IsNull() {
				in.Skip()
				out.Tag = nil
			} else {
				in.Delim('[')
				if out.Tag == nil {
					if !in.IsDelim(']') {
						out.Tag = make([]string, 0, 4)
					} else {
						out.Tag = []string{}
					}
				} else {
					out.Tag = (out.Tag)[:0]
				}
				for !in.IsDelim(']') {
					var v16 string
					v16 = string(in.String())
					out.Tag = append(out.Tag, v16)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "author":
			out.AuthorID = int(in.Int())
		case "event_start":
			out.EventStart = string(in.String())
		case "event_end":
			out.EventEnd = string(in.String())
		case "image":
			out.ImageURL = string(in.String())
		case "Latitude":
			out.Latitude = float64(in.Float64())
		case "Longitude":
			out.Longitude = float64(in.Float64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonF642ad3eEncodeKudagoInternalGatewayEvent7(out *jwriter.Writer, in EventResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int(int(in.ID))
	}
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"location\":"
		out.RawString(prefix)
		out.String(string(in.Location))
	}
	{
		const prefix string = ",\"category_id\":"
		out.RawString(prefix)
		out.Int(int(in.Category))
	}
	{
		const prefix string = ",\"capacity\":"
		out.RawString(prefix)
		out.Int(int(in.Capacity))
	}
	{
		const prefix string = ",\"tag\":"
		out.RawString(prefix)
		if in.Tag == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v17, v18 := range in.Tag {
				if v17 > 0 {
					out.RawByte(',')
				}
				out.String(string(v18))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"author\":"
		out.RawString(prefix)
		out.Int(int(in.AuthorID))
	}
	{
		const prefix string = ",\"event_start\":"
		out.RawString(prefix)
		out.String(string(in.EventStart))
	}
	{
		const prefix string = ",\"event_end\":"
		out.RawString(prefix)
		out.String(string(in.EventEnd))
	}
	{
		const prefix string = ",\"image\":"
		out.RawString(prefix)
		out.String(string(in.ImageURL))
	}
	{
		const prefix string = ",\"Latitude\":"
		out.RawString(prefix)
		out.Float64(float64(in.Latitude))
	}
	{
		const prefix string = ",\"Longitude\":"
		out.RawString(prefix)
		out.Float64(float64(in.Longitude))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v EventResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonF642ad3eEncodeKudagoInternalGatewayEvent7(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v EventResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF642ad3eEncodeKudagoInternalGatewayEvent7(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *EventResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonF642ad3eDecodeKudagoInternalGatewayEvent7(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *EventResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF642ad3eDecodeKudagoInternalGatewayEvent7(l, v)
}
