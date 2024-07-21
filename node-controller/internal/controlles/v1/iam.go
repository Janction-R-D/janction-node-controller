package v1

import (
	"node-controller/common/wrapper"
	"node-controller/internal/proto"
	"node-controller/internal/services/iam"
)

func Login(ctx *wrapper.Context) {
	wrapper.ApiWrapper(ctx, iam.GetLoginService().Login, true, &proto.LoginReq{}, &wrapper.ApiConfig{ReqType: wrapper.CHECKTYPE_JSON})
}

func GetNonce(ctx *wrapper.Context) {
	wrapper.ApiWrapper(ctx, iam.GetLoginService().GetNonce, false, nil, nil)
}

func Logout(ctx *wrapper.Context) {
	wrapper.ApiWrapper(ctx, iam.GetLoginService().Logout, false, nil, nil)
}
