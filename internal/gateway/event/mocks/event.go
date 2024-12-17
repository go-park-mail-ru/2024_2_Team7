// Code generated by MockGen. DO NOT EDIT.
// Source: ../../event/api/event_grpc.pb.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	event "kudago/internal/event/api"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockEventServiceClient is a mock of EventServiceClient interface.
type MockEventServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockEventServiceClientMockRecorder
}

// MockEventServiceClientMockRecorder is the mock recorder for MockEventServiceClient.
type MockEventServiceClientMockRecorder struct {
	mock *MockEventServiceClient
}

// NewMockEventServiceClient creates a new mock instance.
func NewMockEventServiceClient(ctrl *gomock.Controller) *MockEventServiceClient {
	mock := &MockEventServiceClient{ctrl: ctrl}
	mock.recorder = &MockEventServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEventServiceClient) EXPECT() *MockEventServiceClientMockRecorder {
	return m.recorder
}

// AddEvent mocks base method.
func (m *MockEventServiceClient) AddEvent(ctx context.Context, in *event.Event, opts ...grpc.CallOption) (*event.Event, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddEvent", varargs...)
	ret0, _ := ret[0].(*event.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddEvent indicates an expected call of AddEvent.
func (mr *MockEventServiceClientMockRecorder) AddEvent(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddEvent", reflect.TypeOf((*MockEventServiceClient)(nil).AddEvent), varargs...)
}

// AddEventToFavorites mocks base method.
func (m *MockEventServiceClient) AddEventToFavorites(ctx context.Context, in *event.FavoriteEvent, opts ...grpc.CallOption) (*event.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddEventToFavorites", varargs...)
	ret0, _ := ret[0].(*event.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddEventToFavorites indicates an expected call of AddEventToFavorites.
func (mr *MockEventServiceClientMockRecorder) AddEventToFavorites(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddEventToFavorites", reflect.TypeOf((*MockEventServiceClient)(nil).AddEventToFavorites), varargs...)
}

// DeleteEvent mocks base method.
func (m *MockEventServiceClient) DeleteEvent(ctx context.Context, in *event.DeleteEventRequest, opts ...grpc.CallOption) (*event.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteEvent", varargs...)
	ret0, _ := ret[0].(*event.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteEvent indicates an expected call of DeleteEvent.
func (mr *MockEventServiceClientMockRecorder) DeleteEvent(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteEvent", reflect.TypeOf((*MockEventServiceClient)(nil).DeleteEvent), varargs...)
}

// DeleteEventFromFavorites mocks base method.
func (m *MockEventServiceClient) DeleteEventFromFavorites(ctx context.Context, in *event.FavoriteEvent, opts ...grpc.CallOption) (*event.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteEventFromFavorites", varargs...)
	ret0, _ := ret[0].(*event.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteEventFromFavorites indicates an expected call of DeleteEventFromFavorites.
func (mr *MockEventServiceClientMockRecorder) DeleteEventFromFavorites(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteEventFromFavorites", reflect.TypeOf((*MockEventServiceClient)(nil).DeleteEventFromFavorites), varargs...)
}

// GetCategories mocks base method.
func (m *MockEventServiceClient) GetCategories(ctx context.Context, in *event.Empty, opts ...grpc.CallOption) (*event.GetCategoriesResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetCategories", varargs...)
	ret0, _ := ret[0].(*event.GetCategoriesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCategories indicates an expected call of GetCategories.
func (mr *MockEventServiceClientMockRecorder) GetCategories(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCategories", reflect.TypeOf((*MockEventServiceClient)(nil).GetCategories), varargs...)
}

// GetEventByID mocks base method.
func (m *MockEventServiceClient) GetEventByID(ctx context.Context, in *event.GetEventByIDRequest, opts ...grpc.CallOption) (*event.Event, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetEventByID", varargs...)
	ret0, _ := ret[0].(*event.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEventByID indicates an expected call of GetEventByID.
func (mr *MockEventServiceClientMockRecorder) GetEventByID(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEventByID", reflect.TypeOf((*MockEventServiceClient)(nil).GetEventByID), varargs...)
}

// GetEventsByCategory mocks base method.
func (m *MockEventServiceClient) GetEventsByCategory(ctx context.Context, in *event.GetEventsByCategoryRequest, opts ...grpc.CallOption) (*event.Events, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetEventsByCategory", varargs...)
	ret0, _ := ret[0].(*event.Events)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEventsByCategory indicates an expected call of GetEventsByCategory.
func (mr *MockEventServiceClientMockRecorder) GetEventsByCategory(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEventsByCategory", reflect.TypeOf((*MockEventServiceClient)(nil).GetEventsByCategory), varargs...)
}

// GetEventsByIDs mocks base method.
func (m *MockEventServiceClient) GetEventsByIDs(ctx context.Context, in *event.GetEventsByIDsRequest, opts ...grpc.CallOption) (*event.Events, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetEventsByIDs", varargs...)
	ret0, _ := ret[0].(*event.Events)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEventsByIDs indicates an expected call of GetEventsByIDs.
func (mr *MockEventServiceClientMockRecorder) GetEventsByIDs(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEventsByIDs", reflect.TypeOf((*MockEventServiceClient)(nil).GetEventsByIDs), varargs...)
}

// GetEventsByUser mocks base method.
func (m *MockEventServiceClient) GetEventsByUser(ctx context.Context, in *event.GetEventsByUserRequest, opts ...grpc.CallOption) (*event.Events, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetEventsByUser", varargs...)
	ret0, _ := ret[0].(*event.Events)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEventsByUser indicates an expected call of GetEventsByUser.
func (mr *MockEventServiceClientMockRecorder) GetEventsByUser(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEventsByUser", reflect.TypeOf((*MockEventServiceClient)(nil).GetEventsByUser), varargs...)
}

// GetFavorites mocks base method.
func (m *MockEventServiceClient) GetFavorites(ctx context.Context, in *event.GetFavoritesRequest, opts ...grpc.CallOption) (*event.Events, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetFavorites", varargs...)
	ret0, _ := ret[0].(*event.Events)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFavorites indicates an expected call of GetFavorites.
func (mr *MockEventServiceClientMockRecorder) GetFavorites(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFavorites", reflect.TypeOf((*MockEventServiceClient)(nil).GetFavorites), varargs...)
}

// GetPastEvents mocks base method.
func (m *MockEventServiceClient) GetPastEvents(ctx context.Context, in *event.PaginationParams, opts ...grpc.CallOption) (*event.Events, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetPastEvents", varargs...)
	ret0, _ := ret[0].(*event.Events)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPastEvents indicates an expected call of GetPastEvents.
func (mr *MockEventServiceClientMockRecorder) GetPastEvents(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPastEvents", reflect.TypeOf((*MockEventServiceClient)(nil).GetPastEvents), varargs...)
}

// GetSubscribersIDs mocks base method.
func (m *MockEventServiceClient) GetSubscribersIDs(ctx context.Context, in *event.GetSubscribersIDsRequest, opts ...grpc.CallOption) (*event.GetUserIDsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetSubscribersIDs", varargs...)
	ret0, _ := ret[0].(*event.GetUserIDsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubscribersIDs indicates an expected call of GetSubscribersIDs.
func (mr *MockEventServiceClientMockRecorder) GetSubscribersIDs(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubscribersIDs", reflect.TypeOf((*MockEventServiceClient)(nil).GetSubscribersIDs), varargs...)
}

// GetSubscriptionsEvents mocks base method.
func (m *MockEventServiceClient) GetSubscriptionsEvents(ctx context.Context, in *event.GetSubscriptionsRequest, opts ...grpc.CallOption) (*event.Events, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetSubscriptionsEvents", varargs...)
	ret0, _ := ret[0].(*event.Events)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubscriptionsEvents indicates an expected call of GetSubscriptionsEvents.
func (mr *MockEventServiceClientMockRecorder) GetSubscriptionsEvents(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubscriptionsEvents", reflect.TypeOf((*MockEventServiceClient)(nil).GetSubscriptionsEvents), varargs...)
}

// GetUpcomingEvents mocks base method.
func (m *MockEventServiceClient) GetUpcomingEvents(ctx context.Context, in *event.PaginationParams, opts ...grpc.CallOption) (*event.Events, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetUpcomingEvents", varargs...)
	ret0, _ := ret[0].(*event.Events)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUpcomingEvents indicates an expected call of GetUpcomingEvents.
func (mr *MockEventServiceClientMockRecorder) GetUpcomingEvents(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUpcomingEvents", reflect.TypeOf((*MockEventServiceClient)(nil).GetUpcomingEvents), varargs...)
}

// GetUserIDsByFavoriteEvent mocks base method.
func (m *MockEventServiceClient) GetUserIDsByFavoriteEvent(ctx context.Context, in *event.GetUserIDsByFavoriteEventRequest, opts ...grpc.CallOption) (*event.GetUserIDsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetUserIDsByFavoriteEvent", varargs...)
	ret0, _ := ret[0].(*event.GetUserIDsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserIDsByFavoriteEvent indicates an expected call of GetUserIDsByFavoriteEvent.
func (mr *MockEventServiceClientMockRecorder) GetUserIDsByFavoriteEvent(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserIDsByFavoriteEvent", reflect.TypeOf((*MockEventServiceClient)(nil).GetUserIDsByFavoriteEvent), varargs...)
}

// SearchEvents mocks base method.
func (m *MockEventServiceClient) SearchEvents(ctx context.Context, in *event.SearchParams, opts ...grpc.CallOption) (*event.Events, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SearchEvents", varargs...)
	ret0, _ := ret[0].(*event.Events)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchEvents indicates an expected call of SearchEvents.
func (mr *MockEventServiceClientMockRecorder) SearchEvents(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchEvents", reflect.TypeOf((*MockEventServiceClient)(nil).SearchEvents), varargs...)
}

// UpdateEvent mocks base method.
func (m *MockEventServiceClient) UpdateEvent(ctx context.Context, in *event.Event, opts ...grpc.CallOption) (*event.Event, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateEvent", varargs...)
	ret0, _ := ret[0].(*event.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateEvent indicates an expected call of UpdateEvent.
func (mr *MockEventServiceClientMockRecorder) UpdateEvent(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEvent", reflect.TypeOf((*MockEventServiceClient)(nil).UpdateEvent), varargs...)
}

// MockEventServiceServer is a mock of EventServiceServer interface.
type MockEventServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockEventServiceServerMockRecorder
}

// MockEventServiceServerMockRecorder is the mock recorder for MockEventServiceServer.
type MockEventServiceServerMockRecorder struct {
	mock *MockEventServiceServer
}

// NewMockEventServiceServer creates a new mock instance.
func NewMockEventServiceServer(ctrl *gomock.Controller) *MockEventServiceServer {
	mock := &MockEventServiceServer{ctrl: ctrl}
	mock.recorder = &MockEventServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEventServiceServer) EXPECT() *MockEventServiceServerMockRecorder {
	return m.recorder
}

// AddEvent mocks base method.
func (m *MockEventServiceServer) AddEvent(arg0 context.Context, arg1 *event.Event) (*event.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddEvent", arg0, arg1)
	ret0, _ := ret[0].(*event.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddEvent indicates an expected call of AddEvent.
func (mr *MockEventServiceServerMockRecorder) AddEvent(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddEvent", reflect.TypeOf((*MockEventServiceServer)(nil).AddEvent), arg0, arg1)
}

// AddEventToFavorites mocks base method.
func (m *MockEventServiceServer) AddEventToFavorites(arg0 context.Context, arg1 *event.FavoriteEvent) (*event.Empty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddEventToFavorites", arg0, arg1)
	ret0, _ := ret[0].(*event.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddEventToFavorites indicates an expected call of AddEventToFavorites.
func (mr *MockEventServiceServerMockRecorder) AddEventToFavorites(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddEventToFavorites", reflect.TypeOf((*MockEventServiceServer)(nil).AddEventToFavorites), arg0, arg1)
}

// DeleteEvent mocks base method.
func (m *MockEventServiceServer) DeleteEvent(arg0 context.Context, arg1 *event.DeleteEventRequest) (*event.Empty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteEvent", arg0, arg1)
	ret0, _ := ret[0].(*event.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteEvent indicates an expected call of DeleteEvent.
func (mr *MockEventServiceServerMockRecorder) DeleteEvent(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteEvent", reflect.TypeOf((*MockEventServiceServer)(nil).DeleteEvent), arg0, arg1)
}

// DeleteEventFromFavorites mocks base method.
func (m *MockEventServiceServer) DeleteEventFromFavorites(arg0 context.Context, arg1 *event.FavoriteEvent) (*event.Empty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteEventFromFavorites", arg0, arg1)
	ret0, _ := ret[0].(*event.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteEventFromFavorites indicates an expected call of DeleteEventFromFavorites.
func (mr *MockEventServiceServerMockRecorder) DeleteEventFromFavorites(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteEventFromFavorites", reflect.TypeOf((*MockEventServiceServer)(nil).DeleteEventFromFavorites), arg0, arg1)
}

// GetCategories mocks base method.
func (m *MockEventServiceServer) GetCategories(arg0 context.Context, arg1 *event.Empty) (*event.GetCategoriesResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCategories", arg0, arg1)
	ret0, _ := ret[0].(*event.GetCategoriesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCategories indicates an expected call of GetCategories.
func (mr *MockEventServiceServerMockRecorder) GetCategories(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCategories", reflect.TypeOf((*MockEventServiceServer)(nil).GetCategories), arg0, arg1)
}

// GetEventByID mocks base method.
func (m *MockEventServiceServer) GetEventByID(arg0 context.Context, arg1 *event.GetEventByIDRequest) (*event.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEventByID", arg0, arg1)
	ret0, _ := ret[0].(*event.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEventByID indicates an expected call of GetEventByID.
func (mr *MockEventServiceServerMockRecorder) GetEventByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEventByID", reflect.TypeOf((*MockEventServiceServer)(nil).GetEventByID), arg0, arg1)
}

// GetEventsByCategory mocks base method.
func (m *MockEventServiceServer) GetEventsByCategory(arg0 context.Context, arg1 *event.GetEventsByCategoryRequest) (*event.Events, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEventsByCategory", arg0, arg1)
	ret0, _ := ret[0].(*event.Events)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEventsByCategory indicates an expected call of GetEventsByCategory.
func (mr *MockEventServiceServerMockRecorder) GetEventsByCategory(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEventsByCategory", reflect.TypeOf((*MockEventServiceServer)(nil).GetEventsByCategory), arg0, arg1)
}

// GetEventsByIDs mocks base method.
func (m *MockEventServiceServer) GetEventsByIDs(arg0 context.Context, arg1 *event.GetEventsByIDsRequest) (*event.Events, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEventsByIDs", arg0, arg1)
	ret0, _ := ret[0].(*event.Events)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEventsByIDs indicates an expected call of GetEventsByIDs.
func (mr *MockEventServiceServerMockRecorder) GetEventsByIDs(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEventsByIDs", reflect.TypeOf((*MockEventServiceServer)(nil).GetEventsByIDs), arg0, arg1)
}

// GetEventsByUser mocks base method.
func (m *MockEventServiceServer) GetEventsByUser(arg0 context.Context, arg1 *event.GetEventsByUserRequest) (*event.Events, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEventsByUser", arg0, arg1)
	ret0, _ := ret[0].(*event.Events)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEventsByUser indicates an expected call of GetEventsByUser.
func (mr *MockEventServiceServerMockRecorder) GetEventsByUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEventsByUser", reflect.TypeOf((*MockEventServiceServer)(nil).GetEventsByUser), arg0, arg1)
}

// GetFavorites mocks base method.
func (m *MockEventServiceServer) GetFavorites(arg0 context.Context, arg1 *event.GetFavoritesRequest) (*event.Events, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFavorites", arg0, arg1)
	ret0, _ := ret[0].(*event.Events)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFavorites indicates an expected call of GetFavorites.
func (mr *MockEventServiceServerMockRecorder) GetFavorites(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFavorites", reflect.TypeOf((*MockEventServiceServer)(nil).GetFavorites), arg0, arg1)
}

// GetPastEvents mocks base method.
func (m *MockEventServiceServer) GetPastEvents(arg0 context.Context, arg1 *event.PaginationParams) (*event.Events, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPastEvents", arg0, arg1)
	ret0, _ := ret[0].(*event.Events)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPastEvents indicates an expected call of GetPastEvents.
func (mr *MockEventServiceServerMockRecorder) GetPastEvents(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPastEvents", reflect.TypeOf((*MockEventServiceServer)(nil).GetPastEvents), arg0, arg1)
}

// GetSubscribersIDs mocks base method.
func (m *MockEventServiceServer) GetSubscribersIDs(arg0 context.Context, arg1 *event.GetSubscribersIDsRequest) (*event.GetUserIDsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubscribersIDs", arg0, arg1)
	ret0, _ := ret[0].(*event.GetUserIDsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubscribersIDs indicates an expected call of GetSubscribersIDs.
func (mr *MockEventServiceServerMockRecorder) GetSubscribersIDs(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubscribersIDs", reflect.TypeOf((*MockEventServiceServer)(nil).GetSubscribersIDs), arg0, arg1)
}

// GetSubscriptionsEvents mocks base method.
func (m *MockEventServiceServer) GetSubscriptionsEvents(arg0 context.Context, arg1 *event.GetSubscriptionsRequest) (*event.Events, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubscriptionsEvents", arg0, arg1)
	ret0, _ := ret[0].(*event.Events)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubscriptionsEvents indicates an expected call of GetSubscriptionsEvents.
func (mr *MockEventServiceServerMockRecorder) GetSubscriptionsEvents(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubscriptionsEvents", reflect.TypeOf((*MockEventServiceServer)(nil).GetSubscriptionsEvents), arg0, arg1)
}

// GetUpcomingEvents mocks base method.
func (m *MockEventServiceServer) GetUpcomingEvents(arg0 context.Context, arg1 *event.PaginationParams) (*event.Events, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUpcomingEvents", arg0, arg1)
	ret0, _ := ret[0].(*event.Events)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUpcomingEvents indicates an expected call of GetUpcomingEvents.
func (mr *MockEventServiceServerMockRecorder) GetUpcomingEvents(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUpcomingEvents", reflect.TypeOf((*MockEventServiceServer)(nil).GetUpcomingEvents), arg0, arg1)
}

// GetUserIDsByFavoriteEvent mocks base method.
func (m *MockEventServiceServer) GetUserIDsByFavoriteEvent(arg0 context.Context, arg1 *event.GetUserIDsByFavoriteEventRequest) (*event.GetUserIDsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserIDsByFavoriteEvent", arg0, arg1)
	ret0, _ := ret[0].(*event.GetUserIDsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserIDsByFavoriteEvent indicates an expected call of GetUserIDsByFavoriteEvent.
func (mr *MockEventServiceServerMockRecorder) GetUserIDsByFavoriteEvent(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserIDsByFavoriteEvent", reflect.TypeOf((*MockEventServiceServer)(nil).GetUserIDsByFavoriteEvent), arg0, arg1)
}

// SearchEvents mocks base method.
func (m *MockEventServiceServer) SearchEvents(arg0 context.Context, arg1 *event.SearchParams) (*event.Events, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchEvents", arg0, arg1)
	ret0, _ := ret[0].(*event.Events)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchEvents indicates an expected call of SearchEvents.
func (mr *MockEventServiceServerMockRecorder) SearchEvents(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchEvents", reflect.TypeOf((*MockEventServiceServer)(nil).SearchEvents), arg0, arg1)
}

// UpdateEvent mocks base method.
func (m *MockEventServiceServer) UpdateEvent(arg0 context.Context, arg1 *event.Event) (*event.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEvent", arg0, arg1)
	ret0, _ := ret[0].(*event.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateEvent indicates an expected call of UpdateEvent.
func (mr *MockEventServiceServerMockRecorder) UpdateEvent(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEvent", reflect.TypeOf((*MockEventServiceServer)(nil).UpdateEvent), arg0, arg1)
}

// mustEmbedUnimplementedEventServiceServer mocks base method.
func (m *MockEventServiceServer) mustEmbedUnimplementedEventServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedEventServiceServer")
}

// mustEmbedUnimplementedEventServiceServer indicates an expected call of mustEmbedUnimplementedEventServiceServer.
func (mr *MockEventServiceServerMockRecorder) mustEmbedUnimplementedEventServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedEventServiceServer", reflect.TypeOf((*MockEventServiceServer)(nil).mustEmbedUnimplementedEventServiceServer))
}

// MockUnsafeEventServiceServer is a mock of UnsafeEventServiceServer interface.
type MockUnsafeEventServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeEventServiceServerMockRecorder
}

// MockUnsafeEventServiceServerMockRecorder is the mock recorder for MockUnsafeEventServiceServer.
type MockUnsafeEventServiceServerMockRecorder struct {
	mock *MockUnsafeEventServiceServer
}

// NewMockUnsafeEventServiceServer creates a new mock instance.
func NewMockUnsafeEventServiceServer(ctrl *gomock.Controller) *MockUnsafeEventServiceServer {
	mock := &MockUnsafeEventServiceServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeEventServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeEventServiceServer) EXPECT() *MockUnsafeEventServiceServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedEventServiceServer mocks base method.
func (m *MockUnsafeEventServiceServer) mustEmbedUnimplementedEventServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedEventServiceServer")
}

// mustEmbedUnimplementedEventServiceServer indicates an expected call of mustEmbedUnimplementedEventServiceServer.
func (mr *MockUnsafeEventServiceServerMockRecorder) mustEmbedUnimplementedEventServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedEventServiceServer", reflect.TypeOf((*MockUnsafeEventServiceServer)(nil).mustEmbedUnimplementedEventServiceServer))
}