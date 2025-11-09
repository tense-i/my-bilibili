package result

import (
	"net/http"

	"mybilibili/common/xerr"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// ResponseSuccessBean 成功响应结构体
type ResponseSuccessBean struct {
	Code uint32      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// ResponseErrorBean 错误响应结构体
type ResponseErrorBean struct {
	Code uint32 `json:"code"`
	Msg  string `json:"msg"`
}

// HttpResult HTTP 响应封装
func HttpResult(r *http.Request, w http.ResponseWriter, resp interface{}, err error) {
	if err == nil {
		// 成功响应
		httpx.WriteJson(w, http.StatusOK, &ResponseSuccessBean{
			Code: xerr.OK,
			Msg:  xerr.MapErrMsg(xerr.OK),
			Data: resp,
		})
	} else {
		// 错误响应
		errCode := xerr.SERVER_COMMON_ERROR
		errMsg := err.Error()

		// 判断是否为自定义错误
		if e, ok := err.(*xerr.CodeError); ok {
			errCode = e.GetErrCode()
			errMsg = e.GetErrMsg()
		}

		httpx.WriteJson(w, http.StatusOK, &ResponseErrorBean{
			Code: errCode,
			Msg:  errMsg,
		})
	}
}

// Success 成功响应
func Success(data interface{}) *ResponseSuccessBean {
	return &ResponseSuccessBean{
		Code: xerr.OK,
		Msg:  xerr.MapErrMsg(xerr.OK),
		Data: data,
	}
}

// Error 错误响应
func Error(code uint32, msg string) *ResponseErrorBean {
	return &ResponseErrorBean{
		Code: code,
		Msg:  msg,
	}
}

// ParamErrorResult 参数错误响应
func ParamErrorResult(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJson(w, http.StatusBadRequest, &ResponseErrorBean{
		Code: xerr.REQUEST_PARAM_ERROR,
		Msg:  xerr.MapErrMsg(xerr.REQUEST_PARAM_ERROR),
	})
}




