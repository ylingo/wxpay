package wxpay

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type KV struct {
	key string
	val string
}

type sign struct {
}

func newSign() *sign {
	return &sign{}
}

func (s sign) sign(v interface{}, key string, signType SIGNTYPE) string {
	return s.sliceToXml(s.addSign(s.mapToSliceAndSort(s.structToMap(v)), key, signType))
}

func (s sign) checkSign(v interface{}, key string, signType SIGNTYPE, signStr string) (bool, string) {
	var ret bool = false
	newSignStr := s.getSign(s.mapToSliceAndSort(s.structToMap(v)), key, signType)
	if signStr == newSignStr {
		ret = true
	}
	return ret, newSignStr
}

func (s sign) md5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	return strings.ToUpper(hex.EncodeToString(hash.Sum(nil)))
}

func (s sign) hmac_sha256(str string, secret string) string {
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(str))
	return strings.ToUpper(hex.EncodeToString(hash.Sum(nil)))
}

//
func (s sign) addSign(kvs []KV, key string, signType SIGNTYPE) []KV {
	return append(kvs, KV{"sign", s.getSign(kvs, key, signType)})
}

func (s sign) getSign(kvs []KV, key string, signType SIGNTYPE) string {
	var signStr string = ""
	var signBuf bytes.Buffer
	for _, kv := range kvs {
		signBuf.WriteString(fmt.Sprintf("%s=%s&", kv.key, kv.val))
	}
	signBuf.WriteString(fmt.Sprintf("key=%s", key))
	fmt.Println("aaa:" + signBuf.String())
	if HMAC_SHA256 == signType {
		signStr = s.hmac_sha256(signBuf.String(), key)
	} else {
		signStr = s.md5(signBuf.String())
	}
	return signStr
}

//
func (s sign) mapToSliceAndSort(m map[string]string) []KV {
	var kvs []KV
	for k, v := range m {
		kvs = append(kvs, KV{k, v})
	}
	sort.Slice(kvs, func(i, j int) bool {
		return kvs[i].key < kvs[j].key //升序
	})
	return kvs
}

//
func (s sign) sliceToXml(kvs []KV) string {
	var buf bytes.Buffer
	buf.WriteString("<xml>")
	for _, kv := range kvs {
		buf.WriteString(fmt.Sprintf("<%s><![CDATA[%s]]></%s>", kv.key, kv.val, kv.key))
	}
	buf.WriteString("</xml>")
	return buf.String()
}

func (s sign) structToMap(v interface{}) map[string]string {
	vt := reflect.TypeOf(v)
	vv := reflect.ValueOf(v)

	//var m map[string]string =
	m := make(map[string]string, 0)
	for i := 0; i < vt.NumField(); i++ {
		f := vt.Field(i)
		val := vv.FieldByName(f.Name).String()
		key := f.Tag.Get("xml")
		//如果结构体字段值为空，则不添加到MAP
		//名称是"sign"，则不添加MAP
		if val != "" && key != "sign" {
			m[key] = val
		}
	}
	return m
}
