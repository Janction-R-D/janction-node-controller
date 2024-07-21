package proto

import (
	"common/utils"
)

type PageReq struct {
	Page  int `json:"page" form:"page" validate:"required,gte=1"`           // 页码
	Count int `json:"page_size" form:"page_size" validate:"required,gte=1"` // 条数
}

type PageResp struct {
	Count    int         `json:"count"`    // 条数
	Previous interface{} `json:"previous"` // 上一页
	Next     interface{} `json:"next"`     // 下一页
}

type BaseResp struct {
	RequestId string      `json:"request_id"` // 请求ID
	Code      int         `json:"code"`       // 响应码
	Data      interface{} `json:"data"`       // 数据
	Message   string      `json:"message"`    // 面向用户的错误消息，用户可看懂的信息
	Error     string      `json:"error"`      // 面向开发人员的错误消息，用于问题定位和排查
}

func ConvertToPageResp(from *utils.Page) (to *PageResp) {
	return &PageResp{
		Count:    from.Total,
		Previous: from.Pre,
		Next:     from.Next,
	}
}

func GetPageResp(page, count, total int) PageResp {
	p := utils.GetPaginator(page, count, total)
	return PageResp{
		Count:    p.Total,
		Previous: p.Pre,
		Next:     p.Next,
	}
}
