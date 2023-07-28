package xerr

// 错误映射Map[code]错误Msg
// 一般业务错误
var errMsgMap map[uint32]string

// 数据库底层错误
var sysErrMsgMap map[uint32]string

// 错误枚举
const (
	// OK 成功返回
	OK uint32 = 200
	// SERVER_ERROR 系统一般错误
	SERVER_ERROR uint32 = 10010
	// REQUEST_PARAM_ERROR 请求参数错误
	REQUEST_PARAM_ERROR uint32 = 10020
	// TOKEN_EXPIRE_ERROR Token过期
	TOKEN_EXPIRE_ERROR uint32 = 10030

	PB_CHECK_ERR         uint32 = 10041
	PB_LOGIC_CHECK_ERR   uint32 = 10042
	RPC_PAGES_PARAM_ERR  uint32 = 10043
	RPC_INSERT_ERR       uint32 = 10044
	RPC_SEARCH_ERR       uint32 = 10045
	RPC_SEARCH_NOT_FOUND uint32 = 10046
	RPC_DELETE_ERR       uint32 = 10047
	RPC_UPDATE_ERR       uint32 = 10048

	// 数据库错误
	DB_DELETE_ERR uint32 = 10051
	DB_INSERT_ERR uint32 = 10052
	DB_UPDATE_ERR uint32 = 10053
	DB_SEARCH_ERR uint32 = 10054

	// 业务校验错误
	CAPTCHA_GEN_ERR   = 10061
	CAPTCHA_CHECK_ERR = 10062

	FILE_UPLOAD_ERR = 10071

	// AUTH_CHECK_FAILURE 未授权
	AUTH_CHECK_FAILURE = 401
)

func init() {
	errMsgMap = make(map[uint32]string)
	errMsgMap[OK] = "success"
	errMsgMap[SERVER_ERROR] = "服务器开小差啦,稍后再来试一试"
	errMsgMap[REQUEST_PARAM_ERROR] = "参数错误"
	errMsgMap[TOKEN_EXPIRE_ERROR] = "token过期，请重新登陆"
	errMsgMap[PB_CHECK_ERR] = "pb参数校验错误"
	errMsgMap[PB_LOGIC_CHECK_ERR] = "pb逻辑校验错误(查询字段->查询值为空)"
	errMsgMap[RPC_PAGES_PARAM_ERR] = "rpc分页参数校验错误"
	errMsgMap[RPC_INSERT_ERR] = "rpc新建数据错误"
	errMsgMap[RPC_SEARCH_ERR] = "rpc搜索错误"
	errMsgMap[RPC_SEARCH_NOT_FOUND] = "数据不存在"
	errMsgMap[RPC_DELETE_ERR] = "rpc删除错误"
	errMsgMap[RPC_UPDATE_ERR] = "rpc更新错误"
	errMsgMap[CAPTCHA_GEN_ERR] = "图形验证码发送失败"
	errMsgMap[CAPTCHA_CHECK_ERR] = "验证码错误"
	errMsgMap[AUTH_CHECK_FAILURE] = "用户未授权"

	sysErrMsgMap = make(map[uint32]string)
	sysErrMsgMap[DB_DELETE_ERR] = "DB删除失败"
	sysErrMsgMap[DB_INSERT_ERR] = "DB插入失败"
	sysErrMsgMap[DB_UPDATE_ERR] = "DB更新失败"
	sysErrMsgMap[DB_SEARCH_ERR] = "DB查询错误"
}

// GetErrMsg
//
//	@Description: 获取默认错误码信息
//	@param code
//	@return string
func GetErrMsg(code uint32) string {
	if msg, ok := errMsgMap[code]; ok {
		return msg
	} else {
		return "服务器开小差啦,稍后再来试一试"
	}
}

// IsBisCodeErr
//
//	@Description: 是否业务错误码错误
//	@param code
//	@return bool
func IsBisCodeErr(code uint32) bool {
	if _, ok := errMsgMap[code]; ok {
		return true
	} else {
		return false
	}
}

// IsSysCodeErr
//
//	@Description: 是否系统底层错误码错误
//	@param code
//	@return bool
func IsSysCodeErr(code uint32) bool {
	if _, ok := sysErrMsgMap[code]; ok {
		return true
	} else {
		return false
	}
}
