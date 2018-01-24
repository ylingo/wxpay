package wxpay

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type KV struct {
	key string
	val interface{}
}

type sign struct {
}

func newSign() *sign {
	return &sign{}
}

func (s sign) sign(v interface{}, key string, signType SIGNTYPE) string {
	return s.sliceToXml(s.addSign(s.mapToSliceAndSort(s.structToMap(v)), key, signType))
}

func (s sign) checkSign(v map[string]string, key string, signType SIGNTYPE, signStr string) (bool, string) {
	var ret bool = false
	newSignStr := s.getSign(s.mapToSliceAndSort(v), key, signType)
	if signStr == newSignStr {
		ret = true
	}
	return ret, newSignStr
}

func (s sign) md5(str string) string {
	return strings.ToUpper(newEncode().md5(str))
}

func (s sign) hmac_sha256(str string, secret string) string {
	return strings.ToUpper(newEncode().hmac_sha256(str, secret))
}

//对已有字段加签，并将key sing加入KV结构的slice
func (s sign) addSign(kvs []KV, key string, signType SIGNTYPE) []KV {
	return append(kvs, KV{"sign", s.getSign(kvs, key, signType)})
}

func (s sign) getSign(kvs []KV, key string, signType SIGNTYPE) string {
	var signStr string = ""
	var signBuf bytes.Buffer
	for _, kv := range kvs {
		signBuf.WriteString(fmt.Sprintf("%s=%s&", kv.key, kv.val.(string)))
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

//map[string]string转换成KV结构体的slice,并对key进行排序
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

//将slice转换成XML
func (s sign) sliceToXml(kvs []KV) string {
	var buf bytes.Buffer
	buf.WriteString("<xml>")
	for _, kv := range kvs {
		buf.WriteString(fmt.Sprintf("<%s><![CDATA[%s]]></%s>", kv.key, kv.val.(string), kv.key))
	}
	buf.WriteString("</xml>")
	return buf.String()
}

//将结构体转换成map[string]string 字段值均变成string类型
//剔除了字段值为空的字段，如果字段名为sign也需要剔除
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
