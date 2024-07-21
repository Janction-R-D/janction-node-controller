package supports

import (
	"node-controller/common/proto"
	"common/mlog"
	"encoding/csv"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/kataras/iris/v12/context"
	"go.uber.org/zap"
)

var XRequestID = "X-Request-Id"

// http response
func makeResponse(ctx context.Context, data interface{}, msg, errMsg string, statusCode int) {
	_, err := ctx.JSON(proto.BaseResp{
		RequestId: ctx.GetHeader(XRequestID),
		Code:      statusCode,
		Message:   msg,
		Data:      data,
		Error:     errMsg,
	})
	if err != nil {
		mlog.Info("make response error ", zap.Error(err))
	}
}

func makeErrResponse(ctx context.Context, data interface{}, msg, errMsg string, statusCode int) {
	ctx.StatusCode(statusCode)
	_, err := ctx.JSON(proto.BaseResp{
		RequestId: ctx.GetHeader(XRequestID),
		Code:      statusCode,
		Message:   msg,
		Data:      data,
		Error:     errMsg,
	})
	if err != nil {
		mlog.Info("make response error ", zap.Error(err))
	}
}

func makeBytesResponse(ctx context.Context, data []byte) {
	ctx.Header("Content-Type", "video/mp4")
	ctx.Header("Accept-Ranges", "bytes")
	ctx.StatusCode(206)

	_, err := ctx.Write(data)
	if err != nil {
		mlog.Info("make response error ", zap.Error(err))
	}
}

func makeFileResponse(ctx context.Context, data, name string) {
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", url.QueryEscape(name)))

	var err error
	if strings.HasSuffix(data, ".zip") || strings.HasSuffix(data, "tar.gz") {
		err = ctx.ServeFile(data, false)
	} else {
		err = ctx.ServeFile(data, false)
	}
	if err != nil {
		mlog.Info("make response error ", zap.Error(err))
	}
}

func SendApiErrorResponse(ctx context.Context, msg, errMsg string, statusCode int, templateData ...map[string]interface{}) {
	makeResponse(ctx, nil, msg, errMsg, statusCode)
}

func SendApiErrorResponseSetHttpCode(ctx context.Context, msg, errMsg string, statusCode int) {
	makeErrResponse(ctx, nil, msg, errMsg, statusCode)
}

func makeErrResponseSetStatusCode(ctx context.Context, data interface{}, msg, errMsg string, respCode, statusCode int) {
	ctx.StatusCode(statusCode)
	_, err := ctx.JSON(proto.BaseResp{
		RequestId: ctx.GetHeader(XRequestID),
		Code:      respCode,
		Message:   msg,
		Data:      data,
		Error:     errMsg,
	})
	if err != nil {
		mlog.Info("make response error ", zap.Error(err))
	}
}

func SendApiErrorResponseSetHttpCodeRespCode(ctx context.Context, msg, errMsg string, respCode, statusCode int) {
	makeErrResponseSetStatusCode(ctx, nil, msg, errMsg, respCode, statusCode)
}

func SendApiErrorResponseWithData(ctx context.Context, data interface{}, msg, errMsg string, statusCode int) {
	makeResponse(ctx, data, msg, errMsg, statusCode)
}

func SendApiResponse(ctx context.Context, data interface{}) {
	makeResponse(ctx, data, "success", "", 200)
}

func SendApiWithMessageResponse(ctx context.Context, data interface{}, msg string) {
	makeResponse(ctx, data, msg, "", 200)
}

func SendBytesResponse(ctx context.Context, data []byte) {
	makeBytesResponse(ctx, data)
}

func SendFileResponse(ctx context.Context, data, name string) {
	makeFileResponse(ctx, data, name)
}

func SendCSVFileResponse(ctx context.Context, data interface{}, name string) {
	makeCSVFileResponse(ctx, data, name)
}

func SendImageResponse(ctx context.Context, path, filename string) {
	makeImageResponse(ctx, path, filename)
}

func makeImageResponse(ctx context.Context, path, filename string) {
	ctx.ContentType("image/png")
	err := ctx.SendFile(path, filename)
	if err != nil {
		mlog.Infof("fail to read file from %s", path)
	}
}

func makeCSVFileResponse(ctx context.Context, data interface{}, name string) {
	ctx.Header("Content-Type", "text/csv;charset=utf-8")

	name = fmt.Sprintf("%s.csv", name)
	file := filepath.Join(os.TempDir(), name)

	f, err := os.Create(file)
	if err != nil {
		mlog.Infof("create csv file failed: %s", err.Error())
		return
	}
	defer f.Close()
	defer os.Remove(file)

	_, err = f.WriteString("\xEF\xBB\xBF")
	if err != nil {
		mlog.Infof("write csv file failed: %s", err.Error())
		return
	}

	w := csv.NewWriter(f)
	err = w.WriteAll(data.([][]string))
	if err != nil {
		mlog.Infof("write data to csv file failed: %s", err.Error())
		return
	}
	w.Flush()

	err = ctx.SendFile(file, name)
	if err != nil {
		mlog.Infof("send file %s failed: %s", file, err.Error())
	}
}
