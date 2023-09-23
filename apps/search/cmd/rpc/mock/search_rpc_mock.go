// Code generated by MockGen. DO NOT EDIT.
// Source: ./rpc/search/search.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	search "go-zero-douyin/apps/search/cmd/rpc/search"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockSearch is a mock of Search interface.
type MockSearch struct {
	ctrl     *gomock.Controller
	recorder *MockSearchMockRecorder
}

// MockSearchMockRecorder is the mock recorder for MockSearch.
type MockSearchMockRecorder struct {
	mock *MockSearch
}

// NewMockSearch creates a new mock instance.
func NewMockSearch(ctrl *gomock.Controller) *MockSearch {
	mock := &MockSearch{ctrl: ctrl}
	mock.recorder = &MockSearchMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSearch) EXPECT() *MockSearchMockRecorder {
	return m.recorder
}

// CompleteVideo mocks base method.
func (m *MockSearch) CompleteVideo(ctx context.Context, in *search.CompleteVideoReq, opts ...grpc.CallOption) (*search.CompleteVideoResp, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CompleteVideo", varargs...)
	ret0, _ := ret[0].(*search.CompleteVideoResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CompleteVideo indicates an expected call of CompleteVideo.
func (mr *MockSearchMockRecorder) CompleteVideo(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CompleteVideo", reflect.TypeOf((*MockSearch)(nil).CompleteVideo), varargs...)
}

// SearchUser mocks base method.
func (m *MockSearch) SearchUser(ctx context.Context, in *search.SearchUserReq, opts ...grpc.CallOption) (*search.SearchUserResp, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SearchUser", varargs...)
	ret0, _ := ret[0].(*search.SearchUserResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchUser indicates an expected call of SearchUser.
func (mr *MockSearchMockRecorder) SearchUser(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchUser", reflect.TypeOf((*MockSearch)(nil).SearchUser), varargs...)
}

// SearchVideo mocks base method.
func (m *MockSearch) SearchVideo(ctx context.Context, in *search.SearchVideoReq, opts ...grpc.CallOption) (*search.SearchVideoResp, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SearchVideo", varargs...)
	ret0, _ := ret[0].(*search.SearchVideoResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchVideo indicates an expected call of SearchVideo.
func (mr *MockSearchMockRecorder) SearchVideo(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchVideo", reflect.TypeOf((*MockSearch)(nil).SearchVideo), varargs...)
}
