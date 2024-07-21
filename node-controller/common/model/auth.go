package model

import "node-controller/common/id"

type IamUser struct {
	EthAddress  string `json:"ethAddress"`
	Certificate string `json:"certificate"` // 认证
	UserID      int    `json:"user_id"`     // 用户id
	UserName    string `json:"user_name"`   // 用户名
	Role        Role   `json:"role"`        // 角色信息
	NS          string `json:"ns"`          // namespace
}

type AuthResult struct {
	// AccessKey:{AK_id}
	// JWT:{jwt-id}
	EthAddress  string `json:"eth_address"`
	Certificate string `json:"certificate,omitempty"` // 认证
	UserID      int    `json:"user_id,omitempty"`     // 用户id
	UserName    string `json:"user_name,omitempty"`   // 用户名
	Role        Role   `json:"role,omitempty"`        // 角色信息
	NS          string `json:"ns,omitempty"`          // namespace
}

type Role struct {
	ID   int    `json:"id,omitempty"`   // 角色id
	Name string `json:"name,omitempty"` // 角色名
}

type UserInfo struct {
	ID     id.ID
	Name   string
	Mobile string
	Email  string
}
