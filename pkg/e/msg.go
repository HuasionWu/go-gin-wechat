package e

var MsgFlags = map[int]string{
	SUCCESS:                        "ok",
	ERROR:                          "内部错误",
	INVALID_PARAMS:                 "请求参数错误",
	ERROR_EXIST_OBJECT:             "已存在该对象",
	ERROR_NOT_EXIST_OBJECT:         "该对象不存在",
	ERROR_AUTH_CHECK_TOKEN_FAIL:    "Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT: "Token已超时",
	ERROR_AUTH_TOKEN:               "Token生成失败",
	ERROR_AUTH:                     "Token错误",
	ERROR_LOGIN:                    "用户名或密码错误",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
