package wxpay

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type PayOrder struct {
}

func newPayOrder() *PayOrder {
	return &PayOrder{}
}

func (o *PayOrder) order(body, attach, out_trade_no, spbill_create_ip, trade_type string, total_fee int, notifyUrl string) (map[string]string, error) {
	reqEntry := unifiedOrder_Req{
		AppId:          cfg.App_Id,
		MchId:          cfg.Mch_Id,
		NonceStr:       randstr.GetRandString(),
		SignType:       string(cfg.DefaultSignType),
		Body:           body,
		Attach:         attach,
		OutTradeNo:     out_trade_no,
		FeeType:        "CNY",
		TotalFee:       fmt.Sprintf("%d", total_fee),
		SpBillCreateIp: spbill_create_ip,
		NotifyUrl:      notifyUrl,
		TradeType:      "APP",
	}
	signXml := newSign().sign(reqEntry, cfg.Key, cfg.DefaultSignType)
	//fmt.Print(signXml)
	logs("req", signXml)
	webReq := newHttpRequest()
	req, _ := webReq.getRequest("POST", cfg.WXBaseUrl+cfg.UnifiedOrderUrl, signXml)
	req.Header.Add("ContentType", "text/xml")
	var err error
	if response, err := webReq.do(req); err == nil {
		return o.parseReponse(response)
	}
	return nil, err
}

func (o *PayOrder) close(out_trade_no string) (map[string]string, error) {

	reqEntry := closeOrder_Req{
		AppId:      cfg.App_Id,
		MchId:      cfg.Mch_Id,
		NonceStr:   randstr.GetRandString(),
		OutTradeNo: out_trade_no,
	}
	signXml := newSign().sign(reqEntry, cfg.Key, cfg.DefaultSignType)
	//fmt.Print(signXml)
	logs("req", signXml)
	webReq := newHttpRequest()
	req, _ := webReq.getRequest("POST", cfg.WXBaseUrl+cfg.CloseOrderUrl, signXml)
	req.Header.Add("ContentType", "text/xml")
	var err error
	if response, err := webReq.do(req); err == nil {
		return o.parseReponse(response)
	}
	return nil, err
}

func (o *PayOrder) query(out_trade_no string) (map[string]string, error) {
	reqEntry := queryOrder_Req{
		AppId:      cfg.App_Id,
		MchId:      cfg.Mch_Id,
		NonceStr:   randstr.GetRandString(),
		OutTradeNo: out_trade_no,
	}
	signXml := newSign().sign(reqEntry, cfg.Key, cfg.DefaultSignType)
	//fmt.Print(signXml)
	logs("req", signXml)
	webReq := newHttpRequest()
	req, _ := webReq.getRequest("POST", cfg.WXBaseUrl+cfg.OrderQueryUrl, signXml)
	req.Header.Add("ContentType", "text/xml")
	var err error
	if response, err := webReq.do(req); err == nil {
		return o.parseReponse(response)
	}
	return nil, err
}

func (o *PayOrder) payNotify(w http.ResponseWriter, req *http.Request) {
	var respXml string
	if reqMap, err := o.parseRequest(req); err == nil {
		sign_type := reqMap["sign_type"]
		signStr := reqMap["sign"]
		if sign_type == "" {
			sign_type = string(MD5)
		}
		if singRet, _ := newSign().checkSign(reqMap, cfg.Key, SIGNTYPE(sign_type), signStr); singRet {
			if notifyLogic != nil {
				notifyLogic.OnPayNotify(reqMap)
			}
			respXml = `<xml> 
  			<return_code><![CDATA[SUCCESS]]></return_code>
   			<return_msg><![CDATA[OK]]></return_msg>
 			</xml>`
		} else {
			respXml = fmt.Sprintf(`<xml> 
  			<return_code><![CDATA[%s]]></return_code>
   			<return_msg><![CDATA[%s]]></return_msg>
 			</xml>`, "SIGNERROR", "签名错误")
		}
	} else {
		respXml = fmt.Sprintf(`<xml> 
  			<return_code><![CDATA[FAIL]]></return_code>
   			<return_msg><![CDATA[%s]]></return_msg>
 			</xml>`, err.Error())
	}
	logs("notify_resp", respXml)
	w.Write([]byte(respXml))
}

func (o *PayOrder) refund(out_trade_no, out_refund_no string, total_fee, refund_fee int, refund_desc, refund_account string) (map[string]string, error) {
	reqEntry := refund_Req{
		AppId:         cfg.App_Id,
		MchId:         cfg.Mch_Id,
		NonceStr:      randstr.GetRandString(),
		OutTradeNo:    out_trade_no,
		TotalFee:      fmt.Sprintf("%d", total_fee),
		RefundFee:     fmt.Sprintf("%d", refund_fee),
		RefundDesc:    refund_desc,
		RefundAccount: refund_account,
	}
	signXml := newSign().sign(reqEntry, cfg.Key, cfg.DefaultSignType)
	//fmt.Print(signXml)
	logs("req", signXml)
	webReq := newHttpRequest()
	req, _ := webReq.getRequest("POST", cfg.WXBaseUrl+cfg.RefundUrl, signXml)
	req.Header.Add("ContentType", "text/xml")
	var err error
	if response, err := webReq.do(req); err == nil {
		return o.parseReponse(response)
	}
	return nil, err
}

func (o *PayOrder) refundQuery(orderId string, orderIdType ORDERIDTYPE) (map[string]string, error) {
	reqEntry := queryRefund_Req{
		AppId:    cfg.App_Id,
		MchId:    cfg.Mch_Id,
		NonceStr: randstr.GetRandString(),
		OutTradeNo: func() string {
			if OUTTRADENO == orderIdType {
				return orderId
			}
			return ""
		}(),
		TransactionId: func() string {
			if TRANSACTIONID == orderIdType {
				return orderId
			}
			return ""
		}(),
		OutRefundNo: func() string {
			if OUTREFUNDNO == orderIdType {
				return orderId
			}
			return ""
		}(),
		RefundId: func() string {
			if REFUNDID == orderIdType {
				return orderId
			}
			return ""
		}(),
		OffSet: "",
	}
	signXml := newSign().sign(reqEntry, cfg.Key, cfg.DefaultSignType)
	//fmt.Print(signXml)
	logs("req", signXml)
	webReq := newHttpRequest()
	req, _ := webReq.getRequest("POST", cfg.WXBaseUrl+cfg.RefundUrl, signXml)
	req.Header.Add("ContentType", "text/xml")
	var err error
	if response, err := webReq.do(req); err == nil {
		return o.parseReponse(response)
	}
	return nil, err
}

func (o *PayOrder) refundNotify(w http.ResponseWriter, req *http.Request) {
	respXml := ""
	var err error
	if reqMap, err := o.parseRequest(req); err == nil {
		var reqInfo = reqMap["req_info"]
		if reqInfo != "" {
			if _, err := newDecode().aes_256_ecb(reqInfo, cfg.Key); err == nil {
				//decodeInfo
				if notifyLogic != nil {
					notifyLogic.OnRefundNotify(reqMap)
				}
				respXml = `<xml> 
  					<return_code><![CDATA[SUCCESS]]></return_code>
   					<return_msg><![CDATA[OK]]></return_msg>
 					</xml>`
				logs("notify_resp", respXml)
				w.Write([]byte(respXml))
				return
			}
		}
	}
	respXml = fmt.Sprintf(`<xml> 
  			<return_code><![CDATA[FAIL]]></return_code>
   			<return_msg><![CDATA[%s]]></return_msg>
 			</xml>`, err.Error())

	logs("notify_resp", respXml)
	w.Write([]byte(respXml))
}

func (o *PayOrder) parseRequest(request *http.Request) (req map[string]string, err error) {
	defer request.Body.Close()
	var bodyXml []byte
	if bodyXml, err = ioutil.ReadAll(request.Body); err == nil {
		logs("notify", string(bodyXml))
		req = make(map[string]string)
		if err = xml.Unmarshal(bodyXml, (*xmlMap)(&req)); err == nil {
			if strings.ToUpper(req["return_code"]) == "success" {
				if strings.ToUpper(req["result_code"]) == "SUCCESS" ||
					req["result_code"] == "" {
					return
				} else {
					err = fmt.Errorf("%s,%s", req["err_code"], req["err_code_des"])
				}
			} else {
				err = fmt.Errorf("%s,%s", req["return_code"], req["return_msg"])
			}
		}
	}
	return
}

func (o *PayOrder) parseReponse(response *http.Response) (resp map[string]string, err error) {
	defer response.Body.Close()
	var bodyXml []byte
	if bodyXml, err = ioutil.ReadAll(response.Body); err == nil {
		logs("resp", string(bodyXml))
		resp = make(map[string]string)
		if err = xml.Unmarshal(bodyXml, (*xmlMap)(&resp)); err == nil {
			if strings.ToUpper(resp["return_code"]) == "success" {
				if strings.ToUpper(resp["result_code"]) == "SUCCESS" {
					return
				} else {
					err = fmt.Errorf("%s,%s", resp["err_code"], resp["err_code_des"])
				}
			} else {
				err = fmt.Errorf("%s,%s", resp["return_code"], resp["return_msg"])
			}
		}
	}
	return
}
