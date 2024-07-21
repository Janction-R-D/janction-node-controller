package jwts

import (
	"common/utils"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"
	"strconv"
	"time"

	"github.com/kataras/iris/v12/context"
	"github.com/sirupsen/logrus"
)

type OpenApi struct {
	Method     string `json:"method"`
	PostParams []byte `json:"post_params"`
	GetParams  string `json:"get_params"`
	Timestamp  string `json:"timestamp"`
}

func (params OpenApi) Signature(ctx context.Context, ip string) (bool, error) {
	if ctx.GetHeader("Signature") == "" {
		logrus.WithField("hmacList auth", "params don't have Signature").Error("invalid parameters of Signature")
		return false, fmt.Errorf("invalid parameters of Signature")
	}
	if ctx.GetHeader("AccessKey") == "" && ctx.GetHeader("AccessKey") != "HH8hQucpkEJcl1BYLsoN8c6IUvD882UmSMnv3GVitHBT5" {
		logrus.WithField("hmacList auth", "params don't have AccessKey").Error("invalid parameters of AccessKey")
		return false, fmt.Errorf("invalid parameters of AccessKey")
	}
	if ctx.GetHeader("Timestamp") == "" {
		logrus.WithField("hmacList auth", "params don't have timestamp").Error("invalid parameters of timestamp")
		return false, fmt.Errorf("invalid parameters of timestamp")
	}
	if !TimeCheck(ctx.GetHeader("Timestamp")) {
		return false, fmt.Errorf("invalid parameters of timestamp")
	}
	if ctx.GetHeader("SignatureNonce") == "" {
		logrus.WithField("hmacList auth", "SignatureNonce invalid").Error("SignatureNonce invalid")
		return false, fmt.Errorf("invalid parameters of SignatureNonce")
	}

	em := utils.GetExpiringMap()
	if _, res := em.Get(ctx.GetHeader("SignatureNonce")); res {
		logrus.WithField("hmacList auth", "SignatureNonce invalid").Error("SignatureNonce invalid")
		return false, fmt.Errorf("invalid parameters of SignatureNonce")
	}

	// 校验IP是否在白名单里面
	//if memberInfo.AppWhiteIps != "" {
	//	withipList := strings.Split(memberInfo.AppWhiteIps, "\n")
	//	network := utils.IpToNetWork(ip)
	//	if !utils.ExistInList(ip, withipList) && !utils.ExistInList(network, withipList) {
	//		return false, nil, fmt.Errorf("invalid ip addr")
	//	}
	//}

	// TODO 根据ak 获取对应sk
	SecretKey := "SECRET_KEY_VALUE"
	if checkOpenApiSigntureHash(ctx, params.Method, params.GetParams, params.PostParams, SecretKey) {
		return true, nil
	}
	logrus.Warningf("verify hmacList token failed, sinStr: %s, request token: %s, request ip: %s", "", ctx.GetHeader("SignatureNonce"), ip)
	return false, fmt.Errorf("invalid parameters of Signature")
}

type Encoder interface {
	Encode(data []byte) string
}

func (he *hashEncoder) Encode(data []byte) string {
	defer he.h.Reset()
	he.h.Write(data)
	return he.encodeF(he.h.Sum(nil))
}

func Encode(encoder Encoder, data []byte) string {
	return encoder.Encode(data)
}

type hashEncoder struct {
	h       hash.Hash
	encodeF func([]byte) string
}

func hexSha256(data []byte) string {
	return Encode(&hashEncoder{
		h:       sha256.New(),
		encodeF: hex.EncodeToString,
	}, data)
}

func base64HmacSha256(secret, data []byte) string {
	return Encode(&hashEncoder{
		h:       hmac.New(sha256.New, secret),
		encodeF: base64.StdEncoding.EncodeToString,
	}, data)
}

func checkOpenApiSigntureHash(ctx context.Context, method, urlEncoded string, postParams []byte, sk string) bool {
	ts := ctx.GetHeader("Timestamp")
	newSinStr := ts + "\n" + method + "\n" + urlEncoded + "\n" + hexSha256(postParams)
	baseToken := base64HmacSha256([]byte(sk), []byte(newSinStr))
	return baseToken == ctx.GetHeader("Signature")
}

func TimeCheck(timeStampStr string) bool {
	timeNow := time.Now()
	timeStampInt, err := strconv.ParseInt(timeStampStr, 10, 64)
	if err != nil {
		logrus.WithField("hmac auth", err).Errorf("convert timestamp failed, timestamp: %s", timeStampStr)
		return false
	}
	timeStamp := time.Unix(timeStampInt, 0)
	beginTm := timeNow.Add(time.Duration(-10) * time.Minute)
	endTm := timeNow.Add(time.Duration(10) * time.Minute)
	if timeStamp.Before(beginTm) || timeStamp.After(endTm) {
		logrus.WithField("hmac token", err).Errorf("check time is valid failed")
		return false
	}
	return true
}
