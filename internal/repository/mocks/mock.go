// Code generated by MockGen. DO NOT EDIT.
// Source: C:/Users/User/GolandProjects/github.com/Ktuty/gRPC_thumbnails/internal/repository/repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockThumbnail is a mock of Thumbnail interface.
type MockThumbnail struct {
	ctrl     *gomock.Controller
	recorder *MockThumbnailMockRecorder
}

// MockThumbnailMockRecorder is the mock recorder for MockThumbnail.
type MockThumbnailMockRecorder struct {
	mock *MockThumbnail
}

// NewMockThumbnail creates a new mock instance.
func NewMockThumbnail(ctrl *gomock.Controller) *MockThumbnail {
	mock := &MockThumbnail{ctrl: ctrl}
	mock.recorder = &MockThumbnailMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockThumbnail) EXPECT() *MockThumbnailMockRecorder {
	return m.recorder
}

// GetThumbnail mocks base method.
func (m *MockThumbnail) GetThumbnail(videoID string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetThumbnail", videoID)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetThumbnail indicates an expected call of GetThumbnail.
func (mr *MockThumbnailMockRecorder) GetThumbnail(videoID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetThumbnail", reflect.TypeOf((*MockThumbnail)(nil).GetThumbnail), videoID)
}

// SaveThumbnailToRedis mocks base method.
func (m *MockThumbnail) SaveThumbnailToRedis(videoID string, thumbnail []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveThumbnailToRedis", videoID, thumbnail)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveThumbnailToRedis indicates an expected call of SaveThumbnailToRedis.
func (mr *MockThumbnailMockRecorder) SaveThumbnailToRedis(videoID, thumbnail interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveThumbnailToRedis", reflect.TypeOf((*MockThumbnail)(nil).SaveThumbnailToRedis), videoID, thumbnail)
}
