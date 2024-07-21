package model

const (
	TbUser = "tb_user"
)

type User struct {
	ID                 int    `gorm:"column:id;primary_key;auto_increment"`                                                                              // id
	Name               string `gorm:"column:name;type:varchar(255);not null;uniqueIndex:idx_user_ns_identity_name;index:idx_user_ns_active_name_create"` // 用户名
	EthAddress         string `gorm:"column:eth_address" json:"eth_address"`
	FullName           string `gorm:"column:full_name;type:varchar(255)"`          // 全名
	Mobile             string `gorm:"column:mobile;type:varchar(255)"`             // 手机号
	Email              string `gorm:"column:email;type:varchar(255)"`              // 邮箱
	Avatar             string `gorm:"column:avatar;type:longtext"`                 // 头像
	Password           string `gorm:"column:password;type:varchar(255)"`           //密码
	PasswordStrength   int64  `gorm:"column:password_strength"`                    //当前密码强度
	PasswordUpdateTime int64  `gorm:"column:password_update_time"`                 //密码更新时间
	AllowLogin         bool   `gorm:"column:allow_login;type:tinyint(1);not null"` // 是否: 允许登录
	LoginWhiteAddr     string `gorm:"column:login_white_addr;type:longtext"`       // 用户登录白名单
	OTPStatus          int64  `gorm:"column:otp_status"`                           // 跟人otp开关: 0关闭, 1开启, 2强制
	RoleID             uint64 `gorm:"column:role_id;uniqueIndex:idx_user_role"`    // 1:超级管理员 2:普通用户
	ExpireSet          bool   `gorm:"column:expire_set;type:tinyint(1)"`           // 是否: 启用用户过期设置
	ExpireTimestamp    int64  `gorm:"column:expire_timestamp"`                     // 过期时间戳
}

func (u User) TableName() string {
	return TbUser
}

type UserLoginInfo struct {
	ID                 uint64 `gorm:"column:id;primary_key;auto_increment"`
	UserID             uint64 `gorm:"column:user_id;not null;uniqueIndex:idx_user_login_info"`
	LoginType          string `gorm:"column:login_type;type:varchar(50)"`    // 登录方式
	TerminalType       string `gorm:"column:terminal_type;type:varchar(50)"` // 终端类型
	Browser            string `gorm:"column:browser;type:varchar(50)"`       // 浏览器类型
	AccessIP           string `gorm:"column:access_ip;type:varchar(25)"`     // 接入ip
	LastLoginTime      int64  `gorm:"column:last_login_time"`                // 登录时间
	LastLogoutTime     int64  `gorm:"column:last_logout_time"`               // 登出时间
	ExpireTime         int64  `gorm:"column:expire_time"`                    // 登录过期时间戳
	CertificateVersion int64  `gorm:"column:certificate_version"`            // 当前版本登录凭证
	CreatedAt          int64
	UpdatedAt          int64
}
