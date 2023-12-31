// Code generated by MockGen. DO NOT EDIT.
// Source: ./rpc/internal/model/userModel_gen.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	sql "database/sql"
	model "go-zero-douyin/apps/user/cmd/rpc/internal/model"
	reflect "reflect"

	squirrel "github.com/Masterminds/squirrel"
	gomock "github.com/golang/mock/gomock"
	sqlx "github.com/zeromicro/go-zero/core/stores/sqlx"
)

// MockuserModel is a mock of userModel interface.
type MockuserModel struct {
	ctrl     *gomock.Controller
	recorder *MockuserModelMockRecorder
}

// MockuserModelMockRecorder is the mock recorder for MockuserModel.
type MockuserModelMockRecorder struct {
	mock *MockuserModel
}

// NewMockuserModel creates a new mock instance.
func NewMockuserModel(ctrl *gomock.Controller) *MockuserModel {
	mock := &MockuserModel{ctrl: ctrl}
	mock.recorder = &MockuserModelMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockuserModel) EXPECT() *MockuserModelMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockuserModel) Delete(ctx context.Context, session sqlx.Session, id int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, session, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockuserModelMockRecorder) Delete(ctx, session, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockuserModel)(nil).Delete), ctx, session, id)
}

// DeleteSoft mocks base method.
func (m *MockuserModel) DeleteSoft(ctx context.Context, session sqlx.Session, data *model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSoft", ctx, session, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSoft indicates an expected call of DeleteSoft.
func (mr *MockuserModelMockRecorder) DeleteSoft(ctx, session, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSoft", reflect.TypeOf((*MockuserModel)(nil).DeleteSoft), ctx, session, data)
}

// FindAll mocks base method.
func (m *MockuserModel) FindAll(ctx context.Context, rowBuilder squirrel.SelectBuilder, orderBy string) ([]*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", ctx, rowBuilder, orderBy)
	ret0, _ := ret[0].([]*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockuserModelMockRecorder) FindAll(ctx, rowBuilder, orderBy interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockuserModel)(nil).FindAll), ctx, rowBuilder, orderBy)
}

// FindCount mocks base method.
func (m *MockuserModel) FindCount(ctx context.Context, countBuilder squirrel.SelectBuilder, field string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindCount", ctx, countBuilder, field)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindCount indicates an expected call of FindCount.
func (mr *MockuserModelMockRecorder) FindCount(ctx, countBuilder, field interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindCount", reflect.TypeOf((*MockuserModel)(nil).FindCount), ctx, countBuilder, field)
}

// FindOne mocks base method.
func (m *MockuserModel) FindOne(ctx context.Context, id int64) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOne", ctx, id)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOne indicates an expected call of FindOne.
func (mr *MockuserModelMockRecorder) FindOne(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOne", reflect.TypeOf((*MockuserModel)(nil).FindOne), ctx, id)
}

// FindOneByUsernameIsDelete mocks base method.
func (m *MockuserModel) FindOneByUsernameIsDelete(ctx context.Context, username string, isDelete int64) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOneByUsernameIsDelete", ctx, username, isDelete)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOneByUsernameIsDelete indicates an expected call of FindOneByUsernameIsDelete.
func (mr *MockuserModelMockRecorder) FindOneByUsernameIsDelete(ctx, username, isDelete interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOneByUsernameIsDelete", reflect.TypeOf((*MockuserModel)(nil).FindOneByUsernameIsDelete), ctx, username, isDelete)
}

// FindPageListByIdASC mocks base method.
func (m *MockuserModel) FindPageListByIdASC(ctx context.Context, rowBuilder squirrel.SelectBuilder, preMaxId, pageSize int64) ([]*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindPageListByIdASC", ctx, rowBuilder, preMaxId, pageSize)
	ret0, _ := ret[0].([]*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindPageListByIdASC indicates an expected call of FindPageListByIdASC.
func (mr *MockuserModelMockRecorder) FindPageListByIdASC(ctx, rowBuilder, preMaxId, pageSize interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindPageListByIdASC", reflect.TypeOf((*MockuserModel)(nil).FindPageListByIdASC), ctx, rowBuilder, preMaxId, pageSize)
}

// FindPageListByIdDESC mocks base method.
func (m *MockuserModel) FindPageListByIdDESC(ctx context.Context, rowBuilder squirrel.SelectBuilder, preMinId, pageSize int64) ([]*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindPageListByIdDESC", ctx, rowBuilder, preMinId, pageSize)
	ret0, _ := ret[0].([]*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindPageListByIdDESC indicates an expected call of FindPageListByIdDESC.
func (mr *MockuserModelMockRecorder) FindPageListByIdDESC(ctx, rowBuilder, preMinId, pageSize interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindPageListByIdDESC", reflect.TypeOf((*MockuserModel)(nil).FindPageListByIdDESC), ctx, rowBuilder, preMinId, pageSize)
}

// FindPageListByPage mocks base method.
func (m *MockuserModel) FindPageListByPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindPageListByPage", ctx, rowBuilder, page, pageSize, orderBy)
	ret0, _ := ret[0].([]*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindPageListByPage indicates an expected call of FindPageListByPage.
func (mr *MockuserModelMockRecorder) FindPageListByPage(ctx, rowBuilder, page, pageSize, orderBy interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindPageListByPage", reflect.TypeOf((*MockuserModel)(nil).FindPageListByPage), ctx, rowBuilder, page, pageSize, orderBy)
}

// FindPageListByPageWithTotal mocks base method.
func (m *MockuserModel) FindPageListByPageWithTotal(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*model.User, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindPageListByPageWithTotal", ctx, rowBuilder, page, pageSize, orderBy)
	ret0, _ := ret[0].([]*model.User)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// FindPageListByPageWithTotal indicates an expected call of FindPageListByPageWithTotal.
func (mr *MockuserModelMockRecorder) FindPageListByPageWithTotal(ctx, rowBuilder, page, pageSize, orderBy interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindPageListByPageWithTotal", reflect.TypeOf((*MockuserModel)(nil).FindPageListByPageWithTotal), ctx, rowBuilder, page, pageSize, orderBy)
}

// FindSum mocks base method.
func (m *MockuserModel) FindSum(ctx context.Context, sumBuilder squirrel.SelectBuilder, field string) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindSum", ctx, sumBuilder, field)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindSum indicates an expected call of FindSum.
func (mr *MockuserModelMockRecorder) FindSum(ctx, sumBuilder, field interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindSum", reflect.TypeOf((*MockuserModel)(nil).FindSum), ctx, sumBuilder, field)
}

// Insert mocks base method.
func (m *MockuserModel) Insert(ctx context.Context, session sqlx.Session, data *model.User) (sql.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, session, data)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert indicates an expected call of Insert.
func (mr *MockuserModelMockRecorder) Insert(ctx, session, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockuserModel)(nil).Insert), ctx, session, data)
}

// SelectBuilder mocks base method.
func (m *MockuserModel) SelectBuilder() squirrel.SelectBuilder {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectBuilder")
	ret0, _ := ret[0].(squirrel.SelectBuilder)
	return ret0
}

// SelectBuilder indicates an expected call of SelectBuilder.
func (mr *MockuserModelMockRecorder) SelectBuilder() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectBuilder", reflect.TypeOf((*MockuserModel)(nil).SelectBuilder))
}

// Trans mocks base method.
func (m *MockuserModel) Trans(ctx context.Context, fn func(context.Context, sqlx.Session) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trans", ctx, fn)
	ret0, _ := ret[0].(error)
	return ret0
}

// Trans indicates an expected call of Trans.
func (mr *MockuserModelMockRecorder) Trans(ctx, fn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trans", reflect.TypeOf((*MockuserModel)(nil).Trans), ctx, fn)
}

// Update mocks base method.
func (m *MockuserModel) Update(ctx context.Context, session sqlx.Session, data *model.User) (sql.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, session, data)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockuserModelMockRecorder) Update(ctx, session, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockuserModel)(nil).Update), ctx, session, data)
}

// UpdateWithVersion mocks base method.
func (m *MockuserModel) UpdateWithVersion(ctx context.Context, session sqlx.Session, data *model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateWithVersion", ctx, session, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateWithVersion indicates an expected call of UpdateWithVersion.
func (mr *MockuserModelMockRecorder) UpdateWithVersion(ctx, session, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateWithVersion", reflect.TypeOf((*MockuserModel)(nil).UpdateWithVersion), ctx, session, data)
}
