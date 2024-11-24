// Code generated by MockGen. DO NOT EDIT.
// Source: ./auth.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	models "kudago/internal/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAuthService is a mock of AuthService interface.
type MockAuthService struct {
	ctrl     *gomock.Controller
	recorder *MockAuthServiceMockRecorder
}

// MockAuthServiceMockRecorder is the mock recorder for MockAuthService.
type MockAuthServiceMockRecorder struct {
	mock *MockAuthService
}

// NewMockAuthService creates a new mock instance.
func NewMockAuthService(ctrl *gomock.Controller) *MockAuthService {
	mock := &MockAuthService{ctrl: ctrl}
	mock.recorder = &MockAuthServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthService) EXPECT() *MockAuthServiceMockRecorder {
	return m.recorder
}

// CheckCredentials mocks base method.
func (m *MockAuthService) CheckCredentials(ctx context.Context, creds models.Credentials) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckCredentials", ctx, creds)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckCredentials indicates an expected call of CheckCredentials.
func (mr *MockAuthServiceMockRecorder) CheckCredentials(ctx, creds interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckCredentials", reflect.TypeOf((*MockAuthService)(nil).CheckCredentials), ctx, creds)
}

// CheckSession mocks base method.
func (m *MockAuthService) CheckSession(ctx context.Context, cookie string) (models.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckSession", ctx, cookie)
	ret0, _ := ret[0].(models.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckSession indicates an expected call of CheckSession.
func (mr *MockAuthServiceMockRecorder) CheckSession(ctx, cookie interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckSession", reflect.TypeOf((*MockAuthService)(nil).CheckSession), ctx, cookie)
}

// CreateSession mocks base method.
func (m *MockAuthService) CreateSession(ctx context.Context, ID int) (models.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", ctx, ID)
	ret0, _ := ret[0].(models.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MockAuthServiceMockRecorder) CreateSession(ctx, ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockAuthService)(nil).CreateSession), ctx, ID)
}

// DeleteSession mocks base method.
func (m *MockAuthService) DeleteSession(ctx context.Context, token string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSession", ctx, token)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSession indicates an expected call of DeleteSession.
func (mr *MockAuthServiceMockRecorder) DeleteSession(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSession", reflect.TypeOf((*MockAuthService)(nil).DeleteSession), ctx, token)
}

// GetSubscriptions mocks base method.
func (m *MockAuthService) GetSubscriptions(ctx context.Context, ID int) ([]models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubscriptions", ctx, ID)
	ret0, _ := ret[0].([]models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubscriptions indicates an expected call of GetSubscriptions.
func (mr *MockAuthServiceMockRecorder) GetSubscriptions(ctx, ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubscriptions", reflect.TypeOf((*MockAuthService)(nil).GetSubscriptions), ctx, ID)
}

// GetUserByID mocks base method.
func (m *MockAuthService) GetUserByID(ctx context.Context, ID int) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", ctx, ID)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockAuthServiceMockRecorder) GetUserByID(ctx, ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockAuthService)(nil).GetUserByID), ctx, ID)
}

// Register mocks base method.
func (m *MockAuthService) Register(ctx context.Context, registerDTO models.NewUserData) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", ctx, registerDTO)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register.
func (mr *MockAuthServiceMockRecorder) Register(ctx, registerDTO interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockAuthService)(nil).Register), ctx, registerDTO)
}

// Subscribe mocks base method.
func (m *MockAuthService) Subscribe(ctx context.Context, subscription models.Subscription) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subscribe", ctx, subscription)
	ret0, _ := ret[0].(error)
	return ret0
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockAuthServiceMockRecorder) Subscribe(ctx, subscription interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockAuthService)(nil).Subscribe), ctx, subscription)
}

// Unsubscribe mocks base method.
func (m *MockAuthService) Unsubscribe(ctx context.Context, subscription models.Subscription) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unsubscribe", ctx, subscription)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unsubscribe indicates an expected call of Unsubscribe.
func (mr *MockAuthServiceMockRecorder) Unsubscribe(ctx, subscription interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unsubscribe", reflect.TypeOf((*MockAuthService)(nil).Unsubscribe), ctx, subscription)
}

// UpdateUser mocks base method.
func (m *MockAuthService) UpdateUser(ctx context.Context, data models.NewUserData) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", ctx, data)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockAuthServiceMockRecorder) UpdateUser(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockAuthService)(nil).UpdateUser), ctx, data)
}
