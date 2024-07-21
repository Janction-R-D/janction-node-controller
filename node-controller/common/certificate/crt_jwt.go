package certificate

import (
	"node-controller/common/entities"
	"node-controller/common/id"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"

	utilTime "common/time"

	"github.com/pkg/errors"
)

type JwtCertificate struct {
	raw string
	Uid string `json:"UserID"`
	Ns  string `json:"UserNS"`
	Ver int64  `json:"Version"`
	jwt.StandardClaims
}

func NewJwtCertificateFromRaw(token string, claims jwt.Claims) error {
	parser := jwt.Parser{}
	_, _, err := parser.ParseUnverified(token, claims)
	if err != nil {
		return err
	}

	return nil
}

func NewJwtCertificateFromRawAndVerify(token string, parser *jwt.Parser, key string, claims jwt.Claims) error {
	if parser == nil {
		parser = new(jwt.Parser)
	}
	tk, err := parser.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		if _, ok := err.(*jwt.ValidationError); ok {
			return err
		}
		return errors.New("certificate invalid")
	}

	if !tk.Valid {
		return errors.New("certificate expired")
	}

	return nil
}

func (j *JwtCertificate) Valid() error {
	return j.StandardClaims.Valid()
}

func (j *JwtCertificate) ID() id.ID {
	return id.FromString(j.StandardClaims.Id)
}

func (j *JwtCertificate) Name() string {
	return j.StandardClaims.Subject
}

func (j *JwtCertificate) Version() int64 {
	return j.Ver
}

func (j *JwtCertificate) Attribute() map[string]string {
	return map[string]string{}
}

func (j *JwtCertificate) ExpireAt() time.Time {
	return time.Unix(j.StandardClaims.ExpiresAt, 0)
}

func (j *JwtCertificate) Raw() string {
	return j.raw
}

func (j *JwtCertificate) Type() entities.CertificateType {
	return entities.CertificateTypeJWT
}

func (j *JwtCertificate) String() string {
	return fmt.Sprintf("%s:%s", j.Type(), j.ID().String())
}

func (j *JwtCertificate) UserId() id.ID {
	return id.FromString(j.Uid)
}

func (j *JwtCertificate) Namespace() string {
	return j.Ns
}

func (j *JwtCertificate) UseSubRule() bool {
	return false
}

// HttpJwtCrtGetter 从http 请求中解析token.
type HttpJwtCrtGetter func(req *http.Request) string

// defaultHttpJwtCrtGetter 默认http token 获取方式.
// 按优先级 从 query / header / cookie 获取.
func defaultHttpJwtCrtGetter(req *http.Request) (tk string) {
	if req == nil {
		return ""
	}
	defer func() {
		if v, _ := url.PathUnescape(tk); v != "" {
			tk = v
		}
		tk = strings.TrimSpace(tk)
	}()

	if tk = req.Header.Get("Authorization"); tk != "" {
		return
	}

	if ck, err := req.Cookie("token"); err == nil {
		tk = ck.Value
	}

	return
}

type JwtCertProvider struct {
	parser           *jwt.Parser
	httpTokenGetter  HttpJwtCrtGetter
	key              string
	Issue            string
	LoginTokenMaxDur time.Duration
	staticProvider   CertProvider
}

type JwtCertProviderOption func(jcp *JwtCertProvider)

func JwtCertProviderWithGetter(getter HttpJwtCrtGetter) JwtCertProviderOption {
	return func(jcp *JwtCertProvider) {
		if getter != nil {
			jcp.httpTokenGetter = getter
		}
	}
}

func JwtCertProviderWithParser(parser *jwt.Parser) JwtCertProviderOption {
	return func(jcp *JwtCertProvider) {
		if parser != nil {
			jcp.parser = parser
		}
	}
}

func JwtCertProviderWithLoginDur(dur time.Duration) JwtCertProviderOption {
	return func(jcp *JwtCertProvider) {
		if dur != 0 {
			jcp.LoginTokenMaxDur = dur
		}
	}
}

func JwtCertProviderWithStaticCertProvider(provider CertProvider) JwtCertProviderOption {
	return func(jcp *JwtCertProvider) {
		jcp.staticProvider = provider
	}
}

func JwtKey(baseKey string, trait string) string {
	return baseKey + trait
}

func NewJwtCertProvider(key string, options ...JwtCertProviderOption) *JwtCertProvider {
	if key == "" {
		key = "xylx!950"
	}

	res := &JwtCertProvider{
		parser:           new(jwt.Parser),
		httpTokenGetter:  defaultHttpJwtCrtGetter,
		key:              key,
		LoginTokenMaxDur: time.Hour * 12,
		Issue:            "cnd_yx",
	}
	for _, op := range options {
		op(res)
	}
	return res
}

func (j *JwtCertProvider) ParseRaw(raw string) (entities.Certificate, error) {
	token := raw
	idx := strings.Index(token, " ")
	if idx >= 0 {
		token = token[idx+1:]
	}
	var res = &JwtCertificate{raw: token}
	if err := NewJwtCertificateFromRawAndVerify(token, j.parser, j.key, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (j *JwtCertProvider) ParseHttp(req *http.Request) (entities.Certificate, error) {
	token := j.httpTokenGetter(req)
	if token == "" {
		return nil, entities.ErrCertificateNotProvided
	}

	if j.staticProvider != nil {
		cert, err := j.staticProvider.ParseRaw(token)
		if err == nil {
			return cert, nil
		}
	}

	return j.ParseRaw(token)
}

func (j *JwtCertProvider) CreateCertificate(options entities.CreateCertificateOptions) (entities.Certificate, error) {
	if options.UserID.IsZero() {
		return nil, errors.Errorf("user_id is zero")
	}
	if options.Name == "" {
		return nil, errors.Errorf("name is empty")
	}

	expireAt := options.ExpireAt
	if options.ExpireAt.IsZero() {
		expireAt = utilTime.Now().Add(j.LoginTokenMaxDur)
	}

	cer := &JwtCertificate{
		Uid: options.UserID.String(),
		Ver: options.Version,
		StandardClaims: jwt.StandardClaims{
			Audience:  "_iam",
			ExpiresAt: expireAt.Unix(),
			Id:        UUID(),
			IssuedAt:  utilTime.Now().Unix(),
			Issuer:    j.Issue,
			NotBefore: utilTime.Now().Add(time.Second * -30).Unix(), // 前推30秒，防止时间不同步。
			Subject:   options.Name,
		},
	}

	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, cer)
	token, err := tk.SignedString([]byte(j.key))
	if err != nil {
		return nil, err
	}

	cer.raw = token

	return cer, nil
}

func (j *JwtCertProvider) SupportTypes() []entities.CertificateType {
	return []entities.CertificateType{
		entities.CertificateTypeJWT,
	}
}
