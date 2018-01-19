package wxpay

import (
	//"encoding/xml"
	"fmt"
	"testing"
)

func Init() {
	cfg = WxPayConfig{
		App_Id:          "wx2421b1c4370ec43b",
		Mch_Id:          "10000100",
		WXBaseUrl:       "https://api.mch.weixin.qq.com",
		UnifiedOrderUrl: "/pay/unifiedorder",
		OrderQueryUrl:   "/pay/orderquery",
		CloseOrderUrl:   "/pay/closeorder",
		RefundUrl:       "/secapi/pay/refund",
		DownloadUrl:     "/pay/downloadbill",
		ReportUrl:       "/payitil/report",
		NotifyUrl:       "http://171.221.223.246",
		Key:             "192006250b4c09247ec02edce69f6a2d",
		DefaultSignType: MD5,
	}
}
func _TestRandStr(t *testing.T) {
	randStr := NewRandStr(20, []RANDTYPE{NUMBER, LOWER, UPPER})
	var strs []string
	for i := 0; i < 10; i++ {
		str := randStr.GetRandString()
		strs = append(strs, str)
		t.Log(str)
	}

	count := len(strs)
	checkResult := true
	for k := 0; k < count; k++ {
		for m := k + 1; m < count; m++ {
			if strs[k] == strs[m] {
				t.Error("有重复的随机字符串")
				checkResult = false
				break
			}
		}
		if !checkResult {
			break
		}
	}

	if checkResult {
		t.Log("测试结束")
	}
}

func _TestMd5(t *testing.T) {
	t.Log(newSign().md5("appid=wxd930ea5d5a258f4f&body=test&device_info=1000&mch_id=10000100&nonce_str=ibuaiVcKdpRxkhJA&key=192006250b4c09247ec02edce69f6a2d"))
}

func _TestSha256(t *testing.T) {
	t.Log(newSign().hmac_sha256("appid=wxd930ea5d5a258f4f&body=test&device_info=1000&mch_id=10000100&nonce_str=ibuaiVcKdpRxkhJA&key=192006250b4c09247ec02edce69f6a2d", "192006250b4c09247ec02edce69f6a2d"))
}

func _TestXml(t *testing.T) {
	//	var xmlStr = `<xml><return_code>aaaa</return_code><return_msg>bbbb</return_msg></xml>`
	//	var resp Resp
	//	err := xml.Unmarshal([]byte(xmlStr), &resp)
	//	if err != nil {
	//		t.Error(err)
	//	} else {
	//		t.Log(resp)
	//	}

	//	var resp Resp = Resp{Return_code: "aaa", Return_msg: "bbbb"}
	//	xmlStr, _ := xml.Marshal(resp)
	//	t.Log(string(xmlStr))
}

func _TestSign(t *testing.T) {
	Init()
	//signType := MD5
	req := UnifiedOrder_Req{
		AppId:    cfg.App_Id,
		MchId:    cfg.Mch_Id,
		NonceStr: "1add1a30ac87aa2db72f57a2375d8fec", //randstr.GetRandString(),
		//SignType:       string(signType),
		Body:           "APP支付测试",
		Attach:         "支付测试",
		OutTradeNo:     "1415659990",
		TotalFee:       fmt.Sprintf("%d", 1),
		SpBillCreateIp: "14.23.150.211",
		NotifyUrl:      "http://wxpay.wxutil.com/pub_v2/pay/notify.v2.php", //cfg.NotifyUrl,
		TradeType:      "APP",
	}
	//	req := UnifiedOrder_Req{
	//		AppId:      "wxd930ea5d5a258f4f",
	//		Body:       "test",
	//		DeviceInfo: "1000",
	//		MchId:      "10000100",
	//		NonceStr:   "ibuaiVcKdpRxkhJA",
	//	}
	signXml := newSign().sign(req, cfg.Key, cfg.DefaultSignType)
	t.Log(signXml)
}

func _TestCheckSign(t *testing.T) {
	//	xmlStr := `<xml>
	//    <return_code><![CDATA[SUCCESS]]></return_code>
	//    <return_msg><![CDATA[OK]]></return_msg>
	//    <appid><![CDATA[wx2421b1c4370ec43b]]></appid>
	//    <mch_id><![CDATA[10000100]]></mch_id>
	//    <nonce_str><![CDATA[IITRi8Iabbblz1Jc]]></nonce_str>
	//    <sign><![CDATA[7921E432F65EB8ED0CE9755F0E86D72F]]></sign>
	//    <result_code><![CDATA[SUCCESS]]></result_code>
	//    <prepay_id><![CDATA[wx201411101639507cbf6ffd8b0779950874]]></prepay_id>
	//    <trade_type><![CDATA[APP]]></trade_type>
	// </xml>`

	//	s := newSign()
	//	var resp UnifiedOrder_Resp = UnifiedOrder_Resp{}
	//	xml.Unmarshal([]byte(xmlStr), &resp)

	//	oldSignStr := s.getSign(s.mapToSliceAndSort(s.structToMap(resp)), cfg.Key, MD5)
	//	resp.Sign = oldSignStr
	//	t.Log(resp.Sign)

	//	result, signStr := newSign().checkSign(resp, cfg.Key, MD5, resp.Sign)
	//	t.Log(result, signStr)
	//	if result {
	//		t.Log("success")
	//	} else {
	//		t.Error("fail")
	//	}
}
