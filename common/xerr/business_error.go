package xerr

import (
	"fmt"
)

// BisErr 自定义业务错误
type BisErr struct {
	errCode uint32
	errMsg  string
}

// GetErrCode 获取错误码
func (e *BisErr) GetErrCode() uint32 {
	return e.errCode
}

// GetErrMsg 获取错误信息
func (e *BisErr) GetErrMsg() string {
	return e.errMsg
}

// Error 实现error接口
func (e *BisErr) Error() string {
	return fmt.Sprintf("ErrCode:%d，ErrMsg:%s", e.errCode, e.errMsg)
}

func NewErrCodeMsg(errCode uint32, errMsg string) *BisErr {
	return &BisErr{errCode: errCode, errMsg: errMsg}
}

// NewErrCode
//
//	@Description: 使用ErrCode 创建 bisErr
//	@param errCode
//	@return *BisErr
func NewErrCode(errCode uint32) *BisErr {
	return &BisErr{errCode: errCode, errMsg: GetErrMsg(errCode)}
}

// NewErrMsg
//
//	@Description: 创建通用 BisErr（SERVER_ERROR）
//	@param errMsg 错误消息
//	@return *BisErr
func NewErrMsg(errMsg string) *BisErr {
	return &BisErr{errCode: SERVER_ERROR, errMsg: errMsg}
}

func NewFileErrMsg(errMsg string) *BisErr {
	return &BisErr{errCode: FILE_UPLOAD_ERR, errMsg: errMsg}
}
