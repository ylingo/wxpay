package wxpay

import (
	"bytes"
	"crypto/aes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

type encode struct {
}

func newEncode() *encode {
	return &encode{}
}

func (e encode) md5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

func (e encode) hmac_sha256(str string, secret string) string {
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

func (e encode) base64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

type decode struct {
}

func newDecode() *decode {
	return &decode{}
}

func (d decode) base64(str string) string {
	decodeBytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return ""
	}
	return string(decodeBytes)
}

func (d decode) aes_256_ecb(str string, secret string) (string, error) {
	myaes := newMyAes(newEncode().md5(secret), 32)
	encryptData, err := myaes.encrypt([]byte(str))
	if err != nil {
		return "", err
	}
	return string(encryptData), nil
}

type myAes struct {
	Key       string
	BlockSize int
}

func newMyAes(key string, blockSize int) *myAes {
	return &myAes{Key: key, BlockSize: blockSize}
}

func (a *myAes) padding(src []byte) []byte {
	paddingCount := aes.BlockSize - len(src)%aes.BlockSize
	if paddingCount == 0 {
		return src
	}
	return append(src, bytes.Repeat([]byte{byte(0)}, paddingCount)...)
}

func (a *myAes) unPadding(src []byte) []byte {
	for i := len(src) - 1; ; i-- {
		if src[i] != 0 {
			return src[:i+1]
		}
	}
	return nil
}

func (a *myAes) encrypt(src []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(a.Key))
	if err != nil {
		return nil, err
	}
	//padding
	src = a.padding(src)
	encryptData := make([]byte, len(src))
	tmpData := make([]byte, a.BlockSize)

	for index := 0; index < len(src); index += a.BlockSize {
		block.Encrypt(tmpData, src[index:index+a.BlockSize])
		copy(encryptData, tmpData)
	}
	return encryptData, nil
}
