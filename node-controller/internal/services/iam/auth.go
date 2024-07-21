package iam

import (
	"context"
	"fmt"
	"github.com/casbin/casbin"
	irisContext "github.com/kataras/iris/v12/context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"node-controller/common/cache"
	"node-controller/common/entities"
	model2 "node-controller/common/model"
	"node-controller/common/supports"
	"node-controller/config"
	model3 "node-controller/dao/model"
	"node-controller/model"
	"time"
)

type AuthServer struct {
	cache    cache.Cache
	IAuth    AuthInterface
	enforcer *casbin.Enforcer
}

var (
	auth *AuthServer
)

// 初始化Cache
func InitAuthService(i AuthInterface, cacheIns cache.Cache, modelPath, policyPath string) {
	enforcer := casbin.NewEnforcer(modelPath, policyPath)
	auth = &AuthServer{
		IAuth:    i,
		cache:    cacheIns,
		enforcer: enforcer,
	}
}

func GetAuthService() *AuthServer {
	return auth
}

type AuthInterface interface {
	GetUserInfoByUserID(userID int) (*model.User, error)
	GetUserBlackToken(token string) (*model3.UserBlackToken, error)
}

func (a AuthServer) Authn(ctx irisContext.Context) {
	path := ctx.Path()
	//method := ctx.Request().Method
	if checkAllowURL(path) {
		ctx.Next()
		return
	}

	certificate, err := config.Jwt.ParseHttp(ctx.Request())
	if err != nil {
		supports.SendApiErrorResponseSetHttpCode(ctx, err.Error(), err.Error(), 401)
		return
	}

	//校验token已经失效
	black, err := a.IAuth.GetUserBlackToken(certificate.Raw())
	if black != nil {
		supports.SendApiErrorResponseSetHttpCode(ctx, "无效凭证", "无效凭证", 401)
		return
	}

	var user = model2.IamUser{}
	err = a.cache.Get(context.Background(), getCertificateKey(certificate), &user)
	if err != nil {
		userInfo, err := a.IAuth.GetUserInfoByUserID(int(certificate.UserId().IntValue()))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				supports.SendApiErrorResponseSetHttpCode(ctx, "User Not Found", err.Error(), 401)
				return
			}
			supports.SendApiErrorResponseSetHttpCode(ctx, "Get UserInfo Error", err.Error(), 401)
			return
		}
		user.UserID = userInfo.ID
		user.UserName = userInfo.Name
		user.Certificate = certificate.Raw()
		user.Role.ID = int(userInfo.RoleID)
		user.EthAddress = userInfo.EthAddress
		if user.Role.ID == 1 {
			user.Role.Name = "admin"
		} else {
			user.Role.Name = "user"
		}

		_ = a.cache.Set(context.Background(), getCertificateKey(certificate), &user, time.Second*5)
	}
	ctx.Values().Set("userInfo", &user)
	ctx.Next()
}

func getCertificateKey(cert entities.Certificate) string {
	return fmt.Sprintf("cert:%s", cert.String())
}

func checkAllowURL(reqPath string) bool {
	var igUrl = []string{"/api/iam/v1/auth"}
	for _, v := range igUrl {
		if reqPath == v {
			return true
		}
	}
	return false
}
