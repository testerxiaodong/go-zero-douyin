package response

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-douyin/common/xerr"
	"google.golang.org/grpc/status"
	"net/http"
)

// R 通用返回对象
type R struct {
	Code uint32      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// NullJson 空结构
type NullJson struct{}

// Success
//
//	@Description: 成功返回结构
//	@param data
//	@return *R
func Success(data interface{}) *R {
	return &R{200, "OK", data}
}

// ErrR 错误返回对象
type ErrR struct {
	Code uint32 `json:"code"`
	Msg  string `json:"msg"`
}

// Error
//
//	@Description: 错误返回结构
//	@param errCode
//	@param errMsg
//	@return *ErrR
func Error(errCode uint32, errMsg string) *ErrR {
	return &ErrR{errCode, errMsg}
}

// ApiResult http返回
func ApiResult(r *http.Request, w http.ResponseWriter, resp interface{}, err error) {

	if err == nil {
		// 成功返回
		r := Success(resp)
		httpx.WriteJson(w, http.StatusOK, r)
	} else {
		// 错误返回
		errCode := xerr.SERVER_ERROR
		errMsg := xerr.GetErrMsg(errCode)

		causeErr := errors.Cause(err)
		fmt.Printf("value: %v type: %T\n", causeErr, causeErr)
		// err类型
		if e, ok := causeErr.(*xerr.BisErr); ok { // 自定义错误类型
			// 自定义CodeError
			errCode = e.GetErrCode()
			errMsg = e.GetErrMsg()
		} else {
			if grpcErr, ok := status.FromError(causeErr); ok {
				// grpc err错误
				grpcCode := uint32(grpcErr.Code())
				// 是否是自定义错误，不是则是系统错误，不进行返回
				if xerr.IsBisCodeErr(grpcCode) {
					errCode = grpcCode
					errMsg = grpcErr.Message()
				}
			}
		}

		logx.WithContext(r.Context()).Errorf("【API接口错误】 : %+v ", err)

		httpx.WriteJson(w, http.StatusBadRequest, Error(errCode, errMsg))
	}
}

// ParamErrorResult http 参数校验错误返回
func ParamErrorResult(r *http.Request, w http.ResponseWriter, err error) {
	errMsg := fmt.Sprintf("%s ,%s", xerr.GetErrMsg(xerr.REQUEST_PARAM_ERROR), err.Error())
	httpx.WriteJson(w, http.StatusBadRequest, Error(xerr.REQUEST_PARAM_ERROR, errMsg))
}
