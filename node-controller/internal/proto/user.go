package proto

import "node-controller/common/proto"

type UserInfo struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Mobile string `json:"mobile"`
	Role   Role   `json:"role"`
	Roles  []Role `json:"roles"`
}

type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type AdminUserInfoRep struct {
	proto.PageReq
	Search string `json:"search"`
}

type AdminUserInfoResp struct {
	proto.PageResp
	List []AdminUserInfo `json:"list"`
}

type AdminUserInfo struct {
	ID             int    `json:"id"`
	UserName       string `json:"user_name"`
	CreateUserName string `json:"create_user_name"`
	CreateTime     int64  `json:"create_time"`
}

type EditAdminPasswd struct {
	ID     int    `json:"id"`
	Passwd string `json:"passwd"`
}

type UserID struct {
	ID int `json:"id"`
}
