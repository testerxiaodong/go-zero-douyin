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

// DeleteUser mocks base method.
func (m *MockSearch) DeleteUser(ctx context.Context, in *search.DeleteUserDocumentReq, opts ...grpc.CallOption) (*search.DeleteUserDocumentResp, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteUser", varargs...)
	ret0, _ := ret[0].(*search.DeleteUserDocumentResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockSearchMockRecorder) DeleteUser(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockSearch)(nil).DeleteUser), varargs...)
}

// DeleteVideo mocks base method.
func (m *MockSearch) DeleteVideo(ctx context.Context, in *search.DeleteVideoDocumentReq, opts ...grpc.CallOption) (*search.DeleteVideoDocumentResp, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteVideo", varargs...)
	ret0, _ := ret[0].(*search.DeleteVideoDocumentResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteVideo indicates an expected call of DeleteVideo.
func (mr *MockSearchMockRecorder) DeleteVideo(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteVideo", reflect.TypeOf((*MockSearch)(nil).DeleteVideo), varargs...)
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

// SyncUserInfo mocks base method.
func (m *MockSearch) SyncUserInfo(ctx context.Context, in *search.SyncUserInfoReq, opts ...grpc.CallOption) (*search.SyncUserInfoResp, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SyncUserInfo", varargs...)
	ret0, _ := ret[0].(*search.SyncUserInfoResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SyncUserInfo indicates an expected call of SyncUserInfo.
func (mr *MockSearchMockRecorder) SyncUserInfo(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SyncUserInfo", reflect.TypeOf((*MockSearch)(nil).SyncUserInfo), varargs...)
}

// SyncVideoInfo mocks base method.
func (m *MockSearch) SyncVideoInfo(ctx context.Context, in *search.SyncVideoInfoReq, opts ...grpc.CallOption) (*search.SyncVideoInfoResp, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SyncVideoInfo", varargs...)
	ret0, _ := ret[0].(*search.SyncVideoInfoResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SyncVideoInfo indicates an expected call of SyncVideoInfo.
func (mr *MockSearchMockRecorder) SyncVideoInfo(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SyncVideoInfo", reflect.TypeOf((*MockSearch)(nil).SyncVideoInfo), varargs...)
}
