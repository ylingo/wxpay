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

func (o *PayOrder) order(body, attach, out_trade_no, spbill_create_ip, trade_type string, total_fee int, notifyUrl string) (map[string]interface{}, error) {
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
	webReq := newHttpRequest()
	req, _ := webReq.getRequest("POST", cfg.WXBaseUrl+cfg.UnifiedOrderUrl, signXml)
	req.Header.Add("ContentType", "text/xml")
	var err error
	if response, err := webReq.do(req); err == nil {
		return o.parseReponse(response)
	}
	return nil, err
}

func (o *PayOrder) close(out_trade_no string) (map[string]interface{}, error) {

	reqEntry := closeOrder_Req{
		AppId:      cfg.App_Id,
		MchId:      cfg.Mch_Id,
		NonceStr:   randstr.GetRandString(),
		OutTradeNo: out_trade_no,
	}
	signXml := newSign().sign(reqEntry, cfg.Key, cfg.DefaultSignType)
	//fmt.Print(signXml)
	webReq := newHttpRequest()
	req, _ := webReq.getRequest("POST", cfg.WXBaseUrl+cfg.CloseOrderUrl, signXml)
	req.Header.Add("ContentType", "text/xml")
	var err error
	if response, err := webReq.do(req); err == nil {
		return o.parseReponse(response)
	}
	return nil, err
}

func (o *PayOrder) query(out_trade_no string) (map[string]interface{}, error) {
	reqEntry := queryOrder_Req{
		AppId:      cfg.App_Id,
		MchId:      cfg.Mch_Id,
		NonceStr:   randstr.GetRandString(),
		OutTradeNo: out_trade_no,
	}
	signXml := newSign().sign(reqEntry, cfg.Key, cfg.DefaultSignType)
	//fmt.Print(signXml)
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
	if reqMap, err := o.parseRequest(req); err == nil {
		if notifyLogic != nil {
			notifyLogic.OnPayNotify(reqMap)
		}
		w.Write([]byte(`<xml> 
  			<return_code><![CDATA[SUCCESS]]></return_code>
   			<return_msg><![CDATA[OK]]></return_msg>
 			</xml>`))
	} else {
		w.Write([]byte(fmt.Sprintf(`<xml> 
  			<return_code><![CDATA[FAIL]]></return_code>
   			<return_msg><![CDATA[%s]]></return_msg>
 			</xml>`, err.Error())))
	}
}

func (o *PayOrder) refund(out_trade_no, out_refund_no string, total_fee, refund_fee int, refund_desc, refund_account string) (map[string]interface{}, error) {
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
	webReq := newHttpRequest()
	req, _ := webReq.getRequest("POST", cfg.WXBaseUrl+cfg.RefundUrl, signXml)
	req.Header.Add("ContentType", "text/xml")
	var err error
	if response, err := webReq.do(req); err == nil {
		return o.parseReponse(response)
	}
	return nil, err
}

func (o *PayOrder) refundQuery(orderId string, orderIdType ORDERIDTYPE) (map[string]interface{}, error) {
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
	if reqMap, err := o.parseRequest(req); err == nil {
		if notifyLogic != nil {
			notifyLogic.OnRefundNotify(reqMap)
		}
		w.Write([]byte(`<xml> 
  			<return_code><![CDATA[SUCCESS]]></return_code>
   			<return_msg><![CDATA[OK]]></return_msg>
 			</xml>`))
	} else {
		w.Write([]byte(fmt.Sprintf(`<xml> 
  			<return_code><![CDATA[FAIL]]></return_code>
   			<return_msg><![CDATA[%s]]></return_msg>
 			</xml>`, err.Error())))
	}
}

func (o *PayOrder) parseRequest(request *http.Request) (req map[string]interface{}, err error) {
	defer request.Body.Close()
	var bodyXml []byte
	if bodyXml, err = ioutil.ReadAll(request.Body); err == nil {
		req = make(map[string]interface{})
		if err = xml.Unmarshal(bodyXml, &req); err == nil {
			if strings.ToUpper(req["return_code"].(string)) == "success" {
				if strings.ToUpper(req["result_code"].(string)) == "SUCCESS" ||
					req["result_code"] == nil {
					return
				} else {
					err = fmt.Errorf("%s,%s", req["err_code"].(string), req["err_code_des"].(string))
				}
			} else {
				err = fmt.Errorf("%s,%s", req["return_code"].(string), req["return_msg"].(string))
			}
		}
	}
	return
}

func (o *PayOrder) parseReponse(response *http.Response) (resp map[string]interface{}, err error) {
	defer response.Body.Close()
	var bodyXml []byte
	if bodyXml, err = ioutil.ReadAll(response.Body); err == nil {
		resp = make(map[string]interface{})
		if err = xml.Unmarshal(bodyXml, &resp); err == nil {
			if strings.ToUpper(resp["return_code"].(string)) == "success" {
				if strings.ToUpper(resp["result_code"].(string)) == "SUCCESS" {
					return
				} else {
					err = fmt.Errorf("%s,%s", resp["err_code"].(string), resp["err_code_des"].(string))
				}
			} else {
				err = fmt.Errorf("%s,%s", resp["return_code"].(string), resp["return_msg"].(string))
			}
		}
	}
	return
}
