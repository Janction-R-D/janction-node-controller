package entities

//go:generate mockgen -package=mocks -destination mocks/mock_$GOFILE . Certificate,Service

import (
	"node-controller/common/id"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type CertificateType string

const (
	CertificateTypeJWT       CertificateType = "JWT"
	CertificateTypeStatic    CertificateType = "static"
	CertificateTypeAkSk      CertificateType = "AccessKey"
	CertificateTypeExchange  CertificateType = "exchange"
	CertificateTypeTemporary CertificateType = "temp"
)

const (
	AttributeKeyRaw    = "Raw"
	AttributeKeyCookie = "Cookie"

	AttributeKeyJwtAud = "jwt_aud"
	AttributeKeyJwtIss = "jwt_iss"
	AttributeKeyJwtSub = "jwt_sub"
	AttributeKeyJwtJti = "jwt_jti"
)

const (
	CertNameLogin         = "login"
	CertNamePasswordCheck = "password"
)

func CertificateString(tp CertificateType, id string) string {
	return fmt.Sprintf("%s:%s", tp, id)
}

func ParseCertificateString(src string) (tp CertificateType, id string) {
	idx := strings.Index(src, ":")
	if idx >= 0 {
		return CertificateType(src[:idx]), src[idx+1:]
	}
	return "", src
}

type Certificate interface {
	ID() id.ID                    // 获取凭证唯一ID
	Type() CertificateType        // 获取凭证类型
	Name() string                 // 凭证名称
	Version() int64               // 凭证版本信息
	Attribute() map[string]string // 凭证的额外字段
	ExpireAt() time.Time          // 凭证过期时间, 零值表示没有过期时间
	Raw() string                  // 获取凭证原始数据，多值使用分号进行分割

	String() string // 凭证打印 format: type:id

	UserId() id.ID     // 获取用户ID
	Namespace() string // 获取用户所属租户
	UseSubRule() bool  // 是否使用了子规则
}

type CreateCertificateOptions struct {
	Name       string
	UserID     id.ID
	Version    int64
	Attributes map[string]string
	ExpireAt   time.Time
	Enabled    bool
	UseSubRule bool
}

func NewCreateCertificateOptions(name string, userID id.ID) CreateCertificateOptions {
	return CreateCertificateOptions{
		Name:       name,
		UserID:     userID,
		Attributes: map[string]string{},
		ExpireAt:   time.Time{},
		Enabled:    true,
	}
}

type CertificateService interface {
	// AddCertBlack 禁用凭证, expireAt 为零表示无限期禁用.
	AddCertBlack(ctx context.Context, uid id.ID, tp CertificateType, id id.ID, expireAt time.Time) error
	// DelCertBlack 启用凭证.
	DelCertBlack(ctx context.Context, uid id.ID, tp CertificateType, id id.ID) error
	// InCertBlack 凭证是否可用.
	InCertBlack(ctx context.Context, uid id.ID, tp CertificateType, id id.ID) (bool, error)
	// ParseCertFromHttp 从http 请求中解析凭证.
	ParseCertFromHttp(ctx context.Context, tp CertificateType, req *http.Request) (Certificate, error)
	// ParseCertFromRaw 从原始信息中解析凭证.
	ParseCertFromRaw(ctx context.Context, tp CertificateType, raw string) (Certificate, error)
	// CreateCert 创建凭证.
	CreateCert(ctx context.Context, tp CertificateType, options CreateCertificateOptions) (Certificate, error)
	// VerifyCert 验证凭证.
	VerifyCert(ctx context.Context, certificate Certificate) error
	// DisableCert 禁用凭证.
	DisableCert(ctx context.Context, certificate Certificate) error
}

var (
	ErrCertificateNotFound       = errors.New("certificate not found")
	ErrCertificateNotProvided    = errors.New("certificate not provided")
	ErrCertificateMalformed      = errors.New("certificate malformed")
	ErrCertificateExpired        = errors.New("certificate expired")
	ErrCertificateRefreshExpired = errors.New("certificate refresh expired")
	ErrCertificateNotValidYet    = errors.New("certificate not valid yet")
	ErrCertificateInvalid        = errors.New("certificate invalid")
	ErrCertificateVersion        = errors.New("certificate version not support")
	ErrCertificateTypeNotSupport = errors.New("certificate type not support")
)
