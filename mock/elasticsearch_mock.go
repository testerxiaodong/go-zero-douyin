// Code generated by MockGen. DO NOT EDIT.
// Source: ./common/elasticService/elasticsearch.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	elastic "github.com/olivere/elastic/v7"
)

// MockElasticService is a mock of ElasticService interface.
type MockElasticService struct {
	ctrl     *gomock.Controller
	recorder *MockElasticServiceMockRecorder
}

// MockElasticServiceMockRecorder is the mock recorder for MockElasticService.
type MockElasticServiceMockRecorder struct {
	mock *MockElasticService
}

// NewMockElasticService creates a new mock instance.
func NewMockElasticService(ctrl *gomock.Controller) *MockElasticService {
	mock := &MockElasticService{ctrl: ctrl}
	mock.recorder = &MockElasticServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockElasticService) EXPECT() *MockElasticServiceMockRecorder {
	return m.recorder
}

// CreateDocument mocks base method.
func (m *MockElasticService) CreateDocument(ctx context.Context, indexName, id string, req interface{}) (*elastic.IndexResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateDocument", ctx, indexName, id, req)
	ret0, _ := ret[0].(*elastic.IndexResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateDocument indicates an expected call of CreateDocument.
func (mr *MockElasticServiceMockRecorder) CreateDocument(ctx, indexName, id, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDocument", reflect.TypeOf((*MockElasticService)(nil).CreateDocument), ctx, indexName, id, req)
}

// SearchByKeyword mocks base method.
func (m *MockElasticService) SearchByKeyword(ctx context.Context, indexName, keyword string) (*elastic.SearchResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchByKeyword", ctx, indexName, keyword)
	ret0, _ := ret[0].(*elastic.SearchResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchByKeyword indicates an expected call of SearchByKeyword.
func (mr *MockElasticServiceMockRecorder) SearchByKeyword(ctx, indexName, keyword interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchByKeyword", reflect.TypeOf((*MockElasticService)(nil).SearchByKeyword), ctx, indexName, keyword)
}
