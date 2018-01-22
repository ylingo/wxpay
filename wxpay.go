package wxpay

import (
	"net/http"
)

type WxPayConfig struct {
	App_Id          string
	Mch_Id          string
	WXBaseUrl       string
	UnifiedOrderUrl string
	OrderQueryUrl   string
	CloseOrderUrl   string
	RefundUrl       string
	DownloadUrl     string
	Key             string
	DefaultSignType SIGNTYPE
}

var randstr *RandStr = NewRandStr(20, []RANDTYPE{NUMBER, UPPER, LOWER})
var logCallBack LogCallBack
var cfg WxPayConfig = WxPayConfig{}
var notifyLogic NotifyLogic

type WxPay struct {
}

func NewWxPay(appId, mchId, BaseUrl, key string, logcb LogCallBack, notifycb NotifyLogic) *WxPay {
	cfg = WxPayConfig{
		App_Id:          appId,
		Mch_Id:          mchId,
		WXBaseUrl:       BaseUrl,
		UnifiedOrderUrl: "/pay/unifiedorder",
		OrderQueryUrl:   "/pay/orderquery",
		CloseOrderUrl:   "/pay/closeorder ",
		RefundUrl:       "/secapi/pay/refund",
		DownloadUrl:     "/pay/downloadbill",
		Key:             key,
		DefaultSignType: MD5,
	}
	logCallBack = logcb
	notifyLogic = notifycb
	return &WxPay{}
}

func (w WxPay) UnifiedOrder(body, attach, out_trade_no, spbill_create_ip, trade_type string, total_fee int, notifyUrl string) (map[string]string, error) {
	return newPayOrder().order(body, attach, out_trade_no, spbill_create_ip, trade_type, total_fee, notifyUrl)
}

func (w WxPay) OrderQuery(out_trade_no string) (map[string]string, error) {
	return newPayOrder().query(out_trade_no)
}

func (w WxPay) CloseOrder(out_trade_no string) (map[string]string, error) {
	return newPayOrder().close(out_trade_no)
}

func (w WxPay) Refund(out_trade_no, out_refund_no string, total_fee, refund_fee int, refund_desc, refund_account string) (map[string]string, error) {
	return newPayOrder().refund(out_refund_no, out_refund_no, total_fee, refund_fee, refund_desc, refund_account)
}

func (w WxPay) RefundQuery(orderId string, orderIdType ORDERIDTYPE) (map[string]string, error) {
	return newPayOrder().refundQuery(orderId, orderIdType)
}

func (w WxPay) DownLoadBill(bill_date string, bill_type BILLTYPE) (orderMap []map[string]string, sumMap map[string]string, err error) {
	return newPayOrder().download(bill_date, bill_type)
}

func (w WxPay) PayNotify(resp http.ResponseWriter, req *http.Request) {
	newPayOrder().payNotify(resp, req)
}

func (w WxPay) RefundNotify(resp http.ResponseWriter, req *http.Request) {
	newPayOrder().refundNotify(resp, req)
}

func logs(logType string, logContent string) {
	if logCallBack != nil {
		logCallBack.OnLog(logType, logContent)
	}
}
