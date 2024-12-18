// Code generated by MockGen. DO NOT EDIT.
// Source: ../../image/api/image_grpc.pb.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	image "kudago/internal/image/api"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockImageServiceClient is a mock of ImageServiceClient interface.
type MockImageServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockImageServiceClientMockRecorder
}

// MockImageServiceClientMockRecorder is the mock recorder for MockImageServiceClient.
type MockImageServiceClientMockRecorder struct {
	mock *MockImageServiceClient
}

// NewMockImageServiceClient creates a new mock instance.
func NewMockImageServiceClient(ctrl *gomock.Controller) *MockImageServiceClient {
	mock := &MockImageServiceClient{ctrl: ctrl}
	mock.recorder = &MockImageServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockImageServiceClient) EXPECT() *MockImageServiceClientMockRecorder {
	return m.recorder
}

// DeleteImage mocks base method.
func (m *MockImageServiceClient) DeleteImage(ctx context.Context, in *image.DeleteRequest, opts ...grpc.CallOption) (*image.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteImage", varargs...)
	ret0, _ := ret[0].(*image.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteImage indicates an expected call of DeleteImage.
func (mr *MockImageServiceClientMockRecorder) DeleteImage(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteImage", reflect.TypeOf((*MockImageServiceClient)(nil).DeleteImage), varargs...)
}

// UploadImage mocks base method.
func (m *MockImageServiceClient) UploadImage(ctx context.Context, in *image.UploadRequest, opts ...grpc.CallOption) (*image.UploadResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UploadImage", varargs...)
	ret0, _ := ret[0].(*image.UploadResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadImage indicates an expected call of UploadImage.
func (mr *MockImageServiceClientMockRecorder) UploadImage(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadImage", reflect.TypeOf((*MockImageServiceClient)(nil).UploadImage), varargs...)
}

// MockImageServiceServer is a mock of ImageServiceServer interface.
type MockImageServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockImageServiceServerMockRecorder
}

// MockImageServiceServerMockRecorder is the mock recorder for MockImageServiceServer.
type MockImageServiceServerMockRecorder struct {
	mock *MockImageServiceServer
}

// NewMockImageServiceServer creates a new mock instance.
func NewMockImageServiceServer(ctrl *gomock.Controller) *MockImageServiceServer {
	mock := &MockImageServiceServer{ctrl: ctrl}
	mock.recorder = &MockImageServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockImageServiceServer) EXPECT() *MockImageServiceServerMockRecorder {
	return m.recorder
}

// DeleteImage mocks base method.
func (m *MockImageServiceServer) DeleteImage(arg0 context.Context, arg1 *image.DeleteRequest) (*image.Empty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteImage", arg0, arg1)
	ret0, _ := ret[0].(*image.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteImage indicates an expected call of DeleteImage.
func (mr *MockImageServiceServerMockRecorder) DeleteImage(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteImage", reflect.TypeOf((*MockImageServiceServer)(nil).DeleteImage), arg0, arg1)
}

// UploadImage mocks base method.
func (m *MockImageServiceServer) UploadImage(arg0 context.Context, arg1 *image.UploadRequest) (*image.UploadResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadImage", arg0, arg1)
	ret0, _ := ret[0].(*image.UploadResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadImage indicates an expected call of UploadImage.
func (mr *MockImageServiceServerMockRecorder) UploadImage(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadImage", reflect.TypeOf((*MockImageServiceServer)(nil).UploadImage), arg0, arg1)
}

// mustEmbedUnimplementedImageServiceServer mocks base method.
func (m *MockImageServiceServer) mustEmbedUnimplementedImageServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedImageServiceServer")
}

// mustEmbedUnimplementedImageServiceServer indicates an expected call of mustEmbedUnimplementedImageServiceServer.
func (mr *MockImageServiceServerMockRecorder) mustEmbedUnimplementedImageServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedImageServiceServer", reflect.TypeOf((*MockImageServiceServer)(nil).mustEmbedUnimplementedImageServiceServer))
}

// MockUnsafeImageServiceServer is a mock of UnsafeImageServiceServer interface.
type MockUnsafeImageServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeImageServiceServerMockRecorder
}

// MockUnsafeImageServiceServerMockRecorder is the mock recorder for MockUnsafeImageServiceServer.
type MockUnsafeImageServiceServerMockRecorder struct {
	mock *MockUnsafeImageServiceServer
}

// NewMockUnsafeImageServiceServer creates a new mock instance.
func NewMockUnsafeImageServiceServer(ctrl *gomock.Controller) *MockUnsafeImageServiceServer {
	mock := &MockUnsafeImageServiceServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeImageServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeImageServiceServer) EXPECT() *MockUnsafeImageServiceServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedImageServiceServer mocks base method.
func (m *MockUnsafeImageServiceServer) mustEmbedUnimplementedImageServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedImageServiceServer")
}

// mustEmbedUnimplementedImageServiceServer indicates an expected call of mustEmbedUnimplementedImageServiceServer.
func (mr *MockUnsafeImageServiceServerMockRecorder) mustEmbedUnimplementedImageServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedImageServiceServer", reflect.TypeOf((*MockUnsafeImageServiceServer)(nil).mustEmbedUnimplementedImageServiceServer))
}