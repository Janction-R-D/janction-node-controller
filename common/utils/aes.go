package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"
)

var (
	// ErrInvalidBlockSize indicates hash blocksize <= 0.
	ErrInvalidBlockSize = errors.New("invalid blocksize")

	// ErrInvalidPKCS7Data indicates bad input to PKCS7 pad or unpad.
	ErrInvalidPKCS7Data = errors.New("invalid PKCS7 data (empty or not padded)")

	// ErrInvalidPKCS7Padding indicates PKCS7 unpad fails to bad input.
	ErrInvalidPKCS7Padding = errors.New("invalid padding on input")
)

const (
	MSDefKey = "5bf074f5d21e9037448d6166bfb41adf"
	MSDefIV  = "41804b963aec0a3df78ee32ad9490a97"
)

// SimpleEncrypt 用默认的 key 和 iv 加密 plainText 并以 base64 编码的形式返回密文
// 等价于运行 echo -n $plainText | openssl enc -aes-128-cbc -a -iv 41804b963aec0a3df78ee32ad9490a97 -K 5bf074f5d21e9037448d6166bfb41adf -nosalt
func SimpleEncrypt(plainText []byte) (string, error) {
	key, err := hex.DecodeString(MSDefKey)
	if err != nil {
		return "", err
	}
	iv, err := hex.DecodeString(MSDefIV)
	if err != nil {
		return "", err
	}

	cipherBytes, err := AES128CBCEncrypt(plainText, key, iv)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(cipherBytes), nil
}

// SimpleDecrypt 用默认的 key 和 iv 解密以 base64 编码的形式保存的密文
// 1.base64 解码, 得到密文
// 2.密文 aes-128-cbc 解密，等到带 padding 的明文
// 3.把 pkcs7 padding 去掉，等到正确的明文
// 等价于 echo "$base64Cipher" | openssl enc -aes-128-cbc -d -a -iv 41804b963aec0a3df78ee32ad9490a97 -K 5bf074f5d21e9037448d6166bfb41adf -nosalt
func SimpleDecrypt(base64Cipher string) (plainText []byte, err error) {
	key, err := hex.DecodeString(MSDefKey)
	if err != nil {
		return nil, err
	}
	iv, err := hex.DecodeString(MSDefIV)
	if err != nil {
		return nil, err
	}

	cipherText, err := base64.StdEncoding.DecodeString(base64Cipher)
	if err != nil {
		return nil, err
	}

	return AES128CBCDecrypt(cipherText, key, iv)
}

// AES128CBCEncrypt 将 plainText 进行 pkcs7 padding，然后对 key, iv 用 0 补位（超长则裁剪）
func AES128CBCEncrypt(plainText, key, iv []byte) (cipherText []byte, err error) {
	plainText, err = Pkcs7Pad(plainText, aes.BlockSize)
	if err != nil {
		return nil, err
	}

	key = ZeroPad(key, aes.BlockSize)
	iv = ZeroPad(iv, aes.BlockSize)

	return AES128CBCEncryptRaw(plainText, key, iv)
}

func AES128CBCDecrypt(cipherText, key, iv []byte) (plainText []byte, err error) {
	key = ZeroPad(key, aes.BlockSize)
	iv = ZeroPad(iv, aes.BlockSize)

	plainText, err = AES128CBCDecryptRaw(cipherText, key, iv)
	if err != nil {
		return
	}
	return Pkcs7Unpad(plainText, aes.BlockSize)
}

func AES128CBCEncryptRaw(plainText, key, iv []byte) (cipherText []byte, err error) {
	defer func() {
		e := recover()
		if e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// panic if iv is not 16 byte, we catch that and recover
	ecbMode := cipher.NewCBCEncrypter(block, iv)

	cipherText = make([]byte, len(plainText))
	ecbMode.CryptBlocks(cipherText, plainText)

	return cipherText, nil
}

func AES128CBCDecryptRaw(cipherText, key, iv []byte) (plainText []byte, err error) {
	defer func() {
		e := recover()
		if e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// panic if iv is not 16 byte, we catch that and recover
	ecbMode := cipher.NewCBCDecrypter(block, iv)

	plainText = make([]byte, len(cipherText))
	ecbMode.CryptBlocks(plainText, cipherText)

	return
}

// Pkcs7Pad right-pads the given byte slice with 1 to n bytes, where
// n is the block size. The size of the result is x times n, where x
// is at least 1.
func Pkcs7Pad(b []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if b == nil || len(b) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	n := blockSize - (len(b) % blockSize)
	pb := make([]byte, len(b)+n)
	copy(pb, b)
	copy(pb[len(b):], bytes.Repeat([]byte{byte(n)}, n))
	return pb, nil
}

// Pkcs7Unpad validates and unpads data from the given bytes slice.
// The returned value will be 1 to n bytes smaller depending on the
// amount of padding, where n is the block size.
func Pkcs7Unpad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if b == nil || len(b) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	if len(b)%blocksize != 0 {
		return nil, ErrInvalidPKCS7Padding
	}
	c := b[len(b)-1]
	n := int(c)
	if n == 0 || n > len(b) {
		return nil, ErrInvalidPKCS7Padding
	}
	for i := 0; i < n; i++ {
		if b[len(b)-n+i] != c {
			return nil, ErrInvalidPKCS7Padding
		}
	}
	return b[:len(b)-n], nil
}

func ZeroPad(b []byte, blockSize int) []byte {
	var res []byte
	diff := blockSize - len(b)
	if diff == 0 {
		res = b
	} else if diff > 0 {
		res = append(b, bytes.Repeat([]byte{0}, diff)...)
	} else if diff < 0 {
		res = b[:blockSize]
	}
	return res
}

func MD5String(s string) string {
	m := md5.New()
	m.Write([]byte(s))
	return hex.EncodeToString(m.Sum(nil))
}
