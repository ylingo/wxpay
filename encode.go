package wxpay

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
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

type decode struct {
}

func newDecode() *decode {
	return &decode{}
}

func (d decode) base64(str string) string {
	return ""
}

func (d decode) aes_256_ecb(str string, secret string) string {
	return ""
}
