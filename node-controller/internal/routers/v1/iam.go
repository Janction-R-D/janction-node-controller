package v1

import (
	"github.com/kataras/iris/v12"
	"node-controller/common/wrapper"
	v1 "node-controller/internal/controlles/v1"
)

func RegisterIamRouters(party iris.Party) {
	party = party.Party("/iam/v1")
	party.Handle(iris.MethodPost, "/login", wrapper.Handler(v1.Login))
	party.Handle(iris.MethodPost, "/logout", wrapper.Handler(v1.Logout))
	party.Handle(iris.MethodGet, "/nonce", wrapper.Handler(v1.GetNonce))
	party.Handle(iris.MethodGet, "/nonce", wrapper.Handler(v1.GetNonce))
}
