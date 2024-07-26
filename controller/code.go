package controller

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeLoginFailed
	CodeServerBusy

	CodeForbidden
	CodeNeedLogin
	CodeInvalidToken
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:      "Success",
	CodeInvalidParam: "Request params error",
	CodeLoginFailed:  "Fail to login",
	CodeServerBusy:   "Service busy",

	CodeForbidden:    "Forbidden",
	CodeNeedLogin:    "Need login",
	CodeInvalidToken: "Invalid token",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
