package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"github.com/pkg/errors"
	mrand "math/rand"
	"sync"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var RSA = Rsa{m: map[string]*rsa.PrivateKey{}, keys: []string{}}

func init() {
	go RSA.Clean(60 * time.Second)
}

var src = mrand.NewSource(time.Now().UnixNano())

func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

type Rsa struct {
	m    map[string]*rsa.PrivateKey
	keys []string
	sync.RWMutex
}

func (r *Rsa) Clean(dur time.Duration) {
	for {
		if len(r.m) > 1000 {
			for _, key := range r.keys {
				r.Lock()
				delete(r.m, key)
				r.keys = r.keys[1:]
				r.Unlock()
				if len(r.m) <= 1000 {
					break
				}
			}
		}
		time.Sleep(dur)
	}
}

func (r *Rsa) GenerateKey() (sessionID string, public *rsa.PublicKey) {
	private, _ := rsa.GenerateKey(rand.Reader, 1024)
	public = &private.PublicKey

	r.Lock()
	defer r.Unlock()
	for {
		sessionID = RandStringBytesMaskImprSrc(8)
		if _, ok := r.m[sessionID]; ok {
			continue
		}
		r.m[sessionID] = private
		r.keys = append(r.keys, sessionID)
		break
	}
	return
}

func (r *Rsa) DeleteUserSessionInfo(sessionID string) {
	r.Lock()
	delete(r.m, sessionID)
	r.Unlock()
}
func (r *Rsa) Decrypt(sessionID string, ciphertext string) (plainttext string, err error) {
	r.RLock()
	defer r.RUnlock()
	if _, ok := r.m[sessionID]; !ok {
		err = errors.New("凭证失效, 请刷新浏览器")
		return
	}

	private := r.m[sessionID]

	bytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return
	}

	bytes, err = rsa.DecryptPKCS1v15(rand.Reader, private, bytes)
	plainttext = string(bytes)
	return
}
