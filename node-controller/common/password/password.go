package password

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const (
	base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
)

var coder = base64.NewEncoding(base64Table)

func Encode(password string) (string, error) {
	encodeByte, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.Wrap(err, "password encode failed")
	}
	return string(encodeByte), nil
}

func Compare(encodePassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encodePassword), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func EncodeWithMd5(password string, salt string) string {
	h := md5.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum([]byte(salt)))
}

func CompareWithMd5(encodePassword, password, salt string) bool {
	if encodePassword == EncodeWithMd5(password, salt) {
		return true
	}
	return false
}

func EncodeWithBase64(password string) string {
	return coder.EncodeToString([]byte(password))
}

func CompareWithBase64(encodePassword, password string) bool {
	if encodePassword == EncodeWithBase64(password) {
		return true
	}
	return false
}
