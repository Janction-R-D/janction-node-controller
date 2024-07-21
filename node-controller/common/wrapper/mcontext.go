package wrapper

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/pkg/errors"
	"reflect"
	"sync"

	"node-controller/common/model"
)

const (
	Operator   = "operator"
	Readonly   = "readonly"
	SuperAdmin = "super_admin"
	User       = "user"
)

// Context is our custom context.
// Let's implement a context which will give us access
type Context struct {
	context.Context
	UserInfo model.AuthResult
	Language string // 语言, ZhCh, EnUs
}

func (ctx *Context) bindReqBody(reqBodyReceiver interface{}, bodyType CheckType) (err error) {
	if reqBodyReceiver == nil {
		return nil
	}

	switch bodyType {
	case CHECKTYPE_FORM:
		err = ctx.ReadForm(reqBodyReceiver)
	case CHECKTYPE_JSON:
		err = ctx.ReadJSON(reqBodyReceiver)
	}
	if err != nil {
		err = errors.Wrapf(err,
			"failed to bind req-body of req(method=%s, path=%s) to struct '%s' by field tag '%s', req-body: %v",
			ctx.Method(), ctx.Path(), reflect.TypeOf(reqBodyReceiver), getReqBodyTag(bodyType), reqBodyReceiver)
	}

	return
}

var contextPool = sync.Pool{New: func() interface{} {
	return &Context{}
}}

func acquire(original iris.Context) *Context {
	ctx := contextPool.Get().(*Context)
	ctx.Context = original // set the context to the original one in order to have access to iris's implementation.'
	userInfo, ok := ctx.Values().Get("userInfo").(*model.IamUser)
	if ok {
		ctx.UserInfo = model.AuthResult{
			EthAddress:  userInfo.EthAddress,
			Certificate: userInfo.Certificate,
			UserID:      userInfo.UserID,
			UserName:    userInfo.UserName,
			Role:        userInfo.Role,
		}
	}
	return ctx
}

func release(ctx *Context) {
	contextPool.Put(ctx)
}

// Handler will convert our handler of func(*Context) to an iris Handler,
// in order to be compatible with the HTTP API.
func Handler(h func(*Context)) iris.Handler {
	return func(original iris.Context) {
		ctx := acquire(original)
		if ctx == nil {
			return
		}
		h(ctx)
		release(ctx)
	}
}
