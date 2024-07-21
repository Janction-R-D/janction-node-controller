package routers

import (
	"github.com/kataras/iris/v12"
	v1 "node-controller/internal/routers/v1"
	"node-controller/internal/services/iam"
)

func Register(app *iris.Application) {
	app.Use(Cors)
	app.Use(iam.GetAuthService().Authn)

	var appRouter = app.Party("/api/")

	common := app.Party("/")
	{
		common.Options("*", func(ctx iris.Context) {
			ctx.Next()
		})
	}

	v1.RegisterIamRouters(appRouter)

}

func Cors(ctx iris.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	if ctx.Request().Method == "OPTIONS" {
		ctx.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,PATCH,OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Accept, Authorization")
		ctx.StatusCode(204)
		return
	}
	ctx.Next()
}
