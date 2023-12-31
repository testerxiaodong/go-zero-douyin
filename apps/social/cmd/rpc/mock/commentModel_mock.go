// Code generated by MockGen. DO NOT EDIT.
// Source: ./rpc/internal/model/commentModel_gen.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	sql "database/sql"
	model "go-zero-douyin/apps/social/cmd/rpc/internal/model"
	reflect "reflect"

	squirrel "github.com/Masterminds/squirrel"
	gomock "github.com/golang/mock/gomock"
	sqlx "github.com/zeromicro/go-zero/core/stores/sqlx"
)

// MockcommentModel is a mock of commentModel interface.
type MockcommentModel struct {
	ctrl     *gomock.Controller
	recorder *MockcommentModelMockRecorder
}

// MockcommentModelMockRecorder is the mock recorder for MockcommentModel.
type MockcommentModelMockRecorder struct {
	mock *MockcommentModel
}

// NewMockcommentModel creates a new mock instance.
func NewMockcommentModel(ctrl *gomock.Controller) *MockcommentModel {
	mock := &MockcommentModel{ctrl: ctrl}
	mock.recorder = &MockcommentModelMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockcommentModel) EXPECT() *MockcommentModelMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockcommentModel) Delete(ctx context.Context, session sqlx.Session, id int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, session, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockcommentModelMockRecorder) Delete(ctx, session, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockcommentModel)(nil).Delete), ctx, session, id)
}

// DeleteSoft mocks base method.
func (m *MockcommentModel) DeleteSoft(ctx context.Context, session sqlx.Session, data *model.Comment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSoft", ctx, session, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSoft indicates an expected call of DeleteSoft.
func (mr *MockcommentModelMockRecorder) DeleteSoft(ctx, session, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSoft", reflect.TypeOf((*MockcommentModel)(nil).DeleteSoft), ctx, session, data)
}

// FindAll mocks base method.
func (m *MockcommentModel) FindAll(ctx context.Context, rowBuilder squirrel.SelectBuilder, orderBy string) ([]*model.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", ctx, rowBuilder, orderBy)
	ret0, _ := ret[0].([]*model.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockcommentModelMockRecorder) FindAll(ctx, rowBuilder, orderBy interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockcommentModel)(nil).FindAll), ctx, rowBuilder, orderBy)
}

// FindCount mocks base method.
func (m *MockcommentModel) FindCount(ctx context.Context, countBuilder squirrel.SelectBuilder, field string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindCount", ctx, countBuilder, field)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindCount indicates an expected call of FindCount.
func (mr *MockcommentModelMockRecorder) FindCount(ctx, countBuilder, field interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindCount", reflect.TypeOf((*MockcommentModel)(nil).FindCount), ctx, countBuilder, field)
}

// FindOne mocks base method.
func (m *MockcommentModel) FindOne(ctx context.Context, id int64) (*model.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOne", ctx, id)
	ret0, _ := ret[0].(*model.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOne indicates an expected call of FindOne.
func (mr *MockcommentModelMockRecorder) FindOne(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOne", reflect.TypeOf((*MockcommentModel)(nil).FindOne), ctx, id)
}

// FindPageListByIdASC mocks base method.
func (m *MockcommentModel) FindPageListByIdASC(ctx context.Context, rowBuilder squirrel.SelectBuilder, preMaxId, pageSize int64) ([]*model.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindPageListByIdASC", ctx, rowBuilder, preMaxId, pageSize)
	ret0, _ := ret[0].([]*model.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindPageListByIdASC indicates an expected call of FindPageListByIdASC.
func (mr *MockcommentModelMockRecorder) FindPageListByIdASC(ctx, rowBuilder, preMaxId, pageSize interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindPageListByIdASC", reflect.TypeOf((*MockcommentModel)(nil).FindPageListByIdASC), ctx, rowBuilder, preMaxId, pageSize)
}

// FindPageListByIdDESC mocks base method.
func (m *MockcommentModel) FindPageListByIdDESC(ctx context.Context, rowBuilder squirrel.SelectBuilder, preMinId, pageSize int64) ([]*model.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindPageListByIdDESC", ctx, rowBuilder, preMinId, pageSize)
	ret0, _ := ret[0].([]*model.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindPageListByIdDESC indicates an expected call of FindPageListByIdDESC.
func (mr *MockcommentModelMockRecorder) FindPageListByIdDESC(ctx, rowBuilder, preMinId, pageSize interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindPageListByIdDESC", reflect.TypeOf((*MockcommentModel)(nil).FindPageListByIdDESC), ctx, rowBuilder, preMinId, pageSize)
}

// FindPageListByPage mocks base method.
func (m *MockcommentModel) FindPageListByPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*model.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindPageListByPage", ctx, rowBuilder, page, pageSize, orderBy)
	ret0, _ := ret[0].([]*model.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindPageListByPage indicates an expected call of FindPageListByPage.
func (mr *MockcommentModelMockRecorder) FindPageListByPage(ctx, rowBuilder, page, pageSize, orderBy interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindPageListByPage", reflect.TypeOf((*MockcommentModel)(nil).FindPageListByPage), ctx, rowBuilder, page, pageSize, orderBy)
}

// FindPageListByPageWithTotal mocks base method.
func (m *MockcommentModel) FindPageListByPageWithTotal(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*model.Comment, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindPageListByPageWithTotal", ctx, rowBuilder, page, pageSize, orderBy)
	ret0, _ := ret[0].([]*model.Comment)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// FindPageListByPageWithTotal indicates an expected call of FindPageListByPageWithTotal.
func (mr *MockcommentModelMockRecorder) FindPageListByPageWithTotal(ctx, rowBuilder, page, pageSize, orderBy interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindPageListByPageWithTotal", reflect.TypeOf((*MockcommentModel)(nil).FindPageListByPageWithTotal), ctx, rowBuilder, page, pageSize, orderBy)
}

// FindSum mocks base method.
func (m *MockcommentModel) FindSum(ctx context.Context, sumBuilder squirrel.SelectBuilder, field string) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindSum", ctx, sumBuilder, field)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindSum indicates an expected call of FindSum.
func (mr *MockcommentModelMockRecorder) FindSum(ctx, sumBuilder, field interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindSum", reflect.TypeOf((*MockcommentModel)(nil).FindSum), ctx, sumBuilder, field)
}

// Insert mocks base method.
func (m *MockcommentModel) Insert(ctx context.Context, session sqlx.Session, data *model.Comment) (sql.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, session, data)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert indicates an expected call of Insert.
func (mr *MockcommentModelMockRecorder) Insert(ctx, session, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockcommentModel)(nil).Insert), ctx, session, data)
}

// SelectBuilder mocks base method.
func (m *MockcommentModel) SelectBuilder() squirrel.SelectBuilder {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectBuilder")
	ret0, _ := ret[0].(squirrel.SelectBuilder)
	return ret0
}

// SelectBuilder indicates an expected call of SelectBuilder.
func (mr *MockcommentModelMockRecorder) SelectBuilder() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectBuilder", reflect.TypeOf((*MockcommentModel)(nil).SelectBuilder))
}

// Trans mocks base method.
func (m *MockcommentModel) Trans(ctx context.Context, fn func(context.Context, sqlx.Session) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trans", ctx, fn)
	ret0, _ := ret[0].(error)
	return ret0
}

// Trans indicates an expected call of Trans.
func (mr *MockcommentModelMockRecorder) Trans(ctx, fn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trans", reflect.TypeOf((*MockcommentModel)(nil).Trans), ctx, fn)
}

// Update mocks base method.
func (m *MockcommentModel) Update(ctx context.Context, session sqlx.Session, data *model.Comment) (sql.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, session, data)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockcommentModelMockRecorder) Update(ctx, session, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockcommentModel)(nil).Update), ctx, session, data)
}

// UpdateWithVersion mocks base method.
func (m *MockcommentModel) UpdateWithVersion(ctx context.Context, session sqlx.Session, data *model.Comment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateWithVersion", ctx, session, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateWithVersion indicates an expected call of UpdateWithVersion.
func (mr *MockcommentModelMockRecorder) UpdateWithVersion(ctx, session, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateWithVersion", reflect.TypeOf((*MockcommentModel)(nil).UpdateWithVersion), ctx, session, data)
}
