package certificate

import (
	"node-controller/common/entities"
	"github.com/google/uuid"
	"net/http"
)

var UUID = uuid.NewString

type CertProvider interface {
	// ParseRaw 从凭证原始数据中解析凭证.
	ParseRaw(raw string) (entities.Certificate, error)
	// ParseHttp 从http 解析凭证.
	ParseHttp(req *http.Request) (entities.Certificate, error)
	// CreateCertificate 创建凭证.
	CreateCertificate(options entities.CreateCertificateOptions) (entities.Certificate, error)
	// SupportTypes provider 支持的凭证类型
	SupportTypes() []entities.CertificateType
}
