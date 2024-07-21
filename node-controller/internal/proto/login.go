package proto

type LoginReq struct {
	UserName   string `json:"username" form:"username"`
	Password   string `json:"password" form:"password"`
	VerifyCode string `json:"verify_code" form:"verify_code"`
	VerifyID   string `json:"verify_id" form:"verify_id"`
}

type LoginResp struct {
	Token string `json:"token"`
}

type VerifyCode struct {
	ID   string `json:"id"`
	Code string `json:"code,omitempty"`
	Img  string `json:"img"`
}

type AddUser struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
