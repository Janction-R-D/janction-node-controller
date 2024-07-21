package wrapper

import (
	"fmt"
	"runtime"

	"github.com/kataras/iris/v12"
	log "github.com/sirupsen/logrus"

	"node-controller/common/supports"
)

func ApiWrapper(ctx *Context, handler ApiHandler, checkReqBody bool, reqBody interface{}, apiConfig *ApiConfig) {
	defer func() {
		if r := recover(); r != nil {
			var buf [10240]byte
			var n = runtime.Stack(buf[:], false)
			var errContent = fmt.Sprintf("panic: %+v\n%s", r, string(buf[:n]))
			var reqInfo = fmt.Sprintf("request: api: %s, req-body: %+v, http-cli: %s", ctx.Path(), reqBody,
				ctx.RemoteAddr())
			log.Errorf("\n%s\n%s", reqInfo, errContent)
			supports.SendApiErrorResponse(ctx.Context, supports.InternalError, "panic in "+reqInfo, iris.StatusInternalServerError)
		}
	}()

	if reqBody != nil {
		var err = ctx.bindReqBody(reqBody, apiConfig.ReqType)
		if err != nil {
			log.Error(err)
			msg := "解析请求参数失败"
			supports.SendApiErrorResponse(ctx, msg, err.Error(), iris.StatusBadRequest)
			return
		}
	}

	if checkReqBody {
		if err := validateReqBody(reqBody, getReqBodyTag(apiConfig.ReqType)); err != nil {
			log.Error(err)
			if _, ok := err.(*supports.UIError); ok {
				supports.SendApiErrorResponse(ctx, err.Error(), "", iris.StatusBadRequest)
			} else {
				msg := "参数校验未通过，请检查请求参数是否正确"
				supports.SendApiErrorResponse(ctx, msg, err.Error(), iris.StatusBadRequest)
			}
			return
		}
	}

	handler(ctx, reqBody)
}

type (
	ApiHandler func(ctx *Context, reqBody interface{})
)

type CheckType = int32

const (
	_ CheckType = iota
	CHECKTYPE_FORM
	CHECKTYPE_JSON
)

func getReqBodyTag(labelType CheckType) string {
	switch labelType {
	default:
		return ""
	case CHECKTYPE_JSON:
		return "json"
	case CHECKTYPE_FORM:
		return "form"
	}
}

type ApiConfig struct {
	ReqType CheckType
}
