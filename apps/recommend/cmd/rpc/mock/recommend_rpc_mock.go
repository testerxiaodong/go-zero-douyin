// Code generated by MockGen. DO NOT EDIT.
// Source: ./rpc/recommend/recommend.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	recommend "go-zero-douyin/apps/recommend/cmd/rpc/recommend"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockRecommend is a mock of Recommend interface.
type MockRecommend struct {
	ctrl     *gomock.Controller
	recorder *MockRecommendMockRecorder
}

// MockRecommendMockRecorder is the mock recorder for MockRecommend.
type MockRecommendMockRecorder struct {
	mock *MockRecommend
}

// NewMockRecommend creates a new mock instance.
func NewMockRecommend(ctrl *gomock.Controller) *MockRecommend {
	mock := &MockRecommend{ctrl: ctrl}
	mock.recorder = &MockRecommendMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRecommend) EXPECT() *MockRecommendMockRecorder {
	return m.recorder
}

// VideoRecommendSection mocks base method.
func (m *MockRecommend) VideoRecommendSection(ctx context.Context, in *recommend.VideoRecommendSectionReq, opts ...grpc.CallOption) (*recommend.VideoRecommendSectionResp, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "VideoRecommendSection", varargs...)
	ret0, _ := ret[0].(*recommend.VideoRecommendSectionResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VideoRecommendSection indicates an expected call of VideoRecommendSection.
func (mr *MockRecommendMockRecorder) VideoRecommendSection(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VideoRecommendSection", reflect.TypeOf((*MockRecommend)(nil).VideoRecommendSection), varargs...)
}
