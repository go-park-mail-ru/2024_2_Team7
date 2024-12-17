// Code generated by MockGen. DO NOT EDIT.
// Source: ../../notification/api/notification_grpc.pb.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	notification "kudago/internal/notification/api"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockNotificationServiceClient is a mock of NotificationServiceClient interface.
type MockNotificationServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockNotificationServiceClientMockRecorder
}

// MockNotificationServiceClientMockRecorder is the mock recorder for MockNotificationServiceClient.
type MockNotificationServiceClientMockRecorder struct {
	mock *MockNotificationServiceClient
}

// NewMockNotificationServiceClient creates a new mock instance.
func NewMockNotificationServiceClient(ctrl *gomock.Controller) *MockNotificationServiceClient {
	mock := &MockNotificationServiceClient{ctrl: ctrl}
	mock.recorder = &MockNotificationServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNotificationServiceClient) EXPECT() *MockNotificationServiceClientMockRecorder {
	return m.recorder
}

// CreateNotifications mocks base method.
func (m *MockNotificationServiceClient) CreateNotifications(ctx context.Context, in *notification.CreateNotificationsRequest, opts ...grpc.CallOption) (*notification.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateNotifications", varargs...)
	ret0, _ := ret[0].(*notification.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateNotifications indicates an expected call of CreateNotifications.
func (mr *MockNotificationServiceClientMockRecorder) CreateNotifications(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNotifications", reflect.TypeOf((*MockNotificationServiceClient)(nil).CreateNotifications), varargs...)
}

// DeleteNotification mocks base method.
func (m *MockNotificationServiceClient) DeleteNotification(ctx context.Context, in *notification.DeleteNotificationRequest, opts ...grpc.CallOption) (*notification.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteNotification", varargs...)
	ret0, _ := ret[0].(*notification.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteNotification indicates an expected call of DeleteNotification.
func (mr *MockNotificationServiceClientMockRecorder) DeleteNotification(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteNotification", reflect.TypeOf((*MockNotificationServiceClient)(nil).DeleteNotification), varargs...)
}

// GetNotifications mocks base method.
func (m *MockNotificationServiceClient) GetNotifications(ctx context.Context, in *notification.GetNotificationsRequest, opts ...grpc.CallOption) (*notification.GetNotificationsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetNotifications", varargs...)
	ret0, _ := ret[0].(*notification.GetNotificationsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNotifications indicates an expected call of GetNotifications.
func (mr *MockNotificationServiceClientMockRecorder) GetNotifications(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNotifications", reflect.TypeOf((*MockNotificationServiceClient)(nil).GetNotifications), varargs...)
}

// MockNotificationServiceServer is a mock of NotificationServiceServer interface.
type MockNotificationServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockNotificationServiceServerMockRecorder
}

// MockNotificationServiceServerMockRecorder is the mock recorder for MockNotificationServiceServer.
type MockNotificationServiceServerMockRecorder struct {
	mock *MockNotificationServiceServer
}

// NewMockNotificationServiceServer creates a new mock instance.
func NewMockNotificationServiceServer(ctrl *gomock.Controller) *MockNotificationServiceServer {
	mock := &MockNotificationServiceServer{ctrl: ctrl}
	mock.recorder = &MockNotificationServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNotificationServiceServer) EXPECT() *MockNotificationServiceServerMockRecorder {
	return m.recorder
}

// CreateNotifications mocks base method.
func (m *MockNotificationServiceServer) CreateNotifications(arg0 context.Context, arg1 *notification.CreateNotificationsRequest) (*notification.Empty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNotifications", arg0, arg1)
	ret0, _ := ret[0].(*notification.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateNotifications indicates an expected call of CreateNotifications.
func (mr *MockNotificationServiceServerMockRecorder) CreateNotifications(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNotifications", reflect.TypeOf((*MockNotificationServiceServer)(nil).CreateNotifications), arg0, arg1)
}

// DeleteNotification mocks base method.
func (m *MockNotificationServiceServer) DeleteNotification(arg0 context.Context, arg1 *notification.DeleteNotificationRequest) (*notification.Empty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteNotification", arg0, arg1)
	ret0, _ := ret[0].(*notification.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteNotification indicates an expected call of DeleteNotification.
func (mr *MockNotificationServiceServerMockRecorder) DeleteNotification(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteNotification", reflect.TypeOf((*MockNotificationServiceServer)(nil).DeleteNotification), arg0, arg1)
}

// GetNotifications mocks base method.
func (m *MockNotificationServiceServer) GetNotifications(arg0 context.Context, arg1 *notification.GetNotificationsRequest) (*notification.GetNotificationsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNotifications", arg0, arg1)
	ret0, _ := ret[0].(*notification.GetNotificationsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNotifications indicates an expected call of GetNotifications.
func (mr *MockNotificationServiceServerMockRecorder) GetNotifications(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNotifications", reflect.TypeOf((*MockNotificationServiceServer)(nil).GetNotifications), arg0, arg1)
}

// mustEmbedUnimplementedNotificationServiceServer mocks base method.
func (m *MockNotificationServiceServer) mustEmbedUnimplementedNotificationServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedNotificationServiceServer")
}

// mustEmbedUnimplementedNotificationServiceServer indicates an expected call of mustEmbedUnimplementedNotificationServiceServer.
func (mr *MockNotificationServiceServerMockRecorder) mustEmbedUnimplementedNotificationServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedNotificationServiceServer", reflect.TypeOf((*MockNotificationServiceServer)(nil).mustEmbedUnimplementedNotificationServiceServer))
}

// MockUnsafeNotificationServiceServer is a mock of UnsafeNotificationServiceServer interface.
type MockUnsafeNotificationServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeNotificationServiceServerMockRecorder
}

// MockUnsafeNotificationServiceServerMockRecorder is the mock recorder for MockUnsafeNotificationServiceServer.
type MockUnsafeNotificationServiceServerMockRecorder struct {
	mock *MockUnsafeNotificationServiceServer
}

// NewMockUnsafeNotificationServiceServer creates a new mock instance.
func NewMockUnsafeNotificationServiceServer(ctrl *gomock.Controller) *MockUnsafeNotificationServiceServer {
	mock := &MockUnsafeNotificationServiceServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeNotificationServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeNotificationServiceServer) EXPECT() *MockUnsafeNotificationServiceServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedNotificationServiceServer mocks base method.
func (m *MockUnsafeNotificationServiceServer) mustEmbedUnimplementedNotificationServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedNotificationServiceServer")
}

// mustEmbedUnimplementedNotificationServiceServer indicates an expected call of mustEmbedUnimplementedNotificationServiceServer.
func (mr *MockUnsafeNotificationServiceServerMockRecorder) mustEmbedUnimplementedNotificationServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedNotificationServiceServer", reflect.TypeOf((*MockUnsafeNotificationServiceServer)(nil).mustEmbedUnimplementedNotificationServiceServer))
}