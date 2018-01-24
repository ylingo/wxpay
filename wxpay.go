package wxpay

import (
	"fmt"
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

//统一下单，该函数能向微信请求下单，并返回prepay_id【微信预支付回话标识】
//body:商品描述
//attach:附加数据，在查询和通知中原样返回
//out_trade_no:商户内部订单号，32字节内，只能是数字和大小写字母
//spbill_create_ip:
//trade_type:交易类型
//total_fee: 订单总金额，单位为分
//nofifyurl:接收微信异步通知的回调地址，不能带参数
func (w WxPay) UnifiedOrder(body, attach, out_trade_no, spbill_create_ip, trade_type string, total_fee int, notifyUrl string) (map[string]string, error) {
	return newPayOrder().order(body, attach, out_trade_no, spbill_create_ip, trade_type, total_fee, notifyUrl)
}

//查询订单
//orderId：在统一下单中定义的商户订单号或者支付结果通知中的微信订单号
//orderIdType: 订单号类型，取值为OUTTRADENO和TRANSACTIONID之一
func (w WxPay) OrderQuery(orderId string, orderIdType ORDERIDTYPE) (map[string]string, error) {
	if orderIdType != TRANSACTIONID &&
		orderIdType != OUTTRADENO {
		return nil, fmt.Errorf("查询订单只能使用商户订单号和微信订单号其中之一")
	}
	return newPayOrder().query(orderId, orderIdType)
}

//关闭订单，当定制支付失败或支付超市或者系统退出等原因，需要调用该接口关闭订单，避免用户重复支付
//out_trade_no：商户订单号
//ps:订单生成后不能马上关闭订单，最短时间为预支付5分钟后
func (w WxPay) CloseOrder(out_trade_no string) (map[string]string, error) {
	return newPayOrder().close(out_trade_no)
}

//申请退款，由于买家或卖家原因需要退款时，卖家调用该接口进行退款
//out_trade_no：商户订单号
//out_refund_no：商户退款单号
//total_fee:订单总金额，单位为分
//refund_fee:退款总金额，单位为分
//refund_desc:退款原因，如果传入非空，则在下发给用户的退款消息中体现退款原因
func (w WxPay) Refund(out_trade_no, out_refund_no string, total_fee, refund_fee int, refund_desc string) (map[string]string, error) {
	return newPayOrder().refund(out_refund_no, out_refund_no, total_fee, refund_fee, refund_desc)
}

//提交退款申请后，通过调用该接口查询退款状态。
//ps:退款有一定延时，用零钱支付的退款20分钟内到账，银行卡支付的退款3个工作日后重新查询退款状态。
//ps:如果单个支付订单部分退款次数超过20次请使用退款单号查询
//orderId:订单ID，包括：TRANSACTIONID\OUTTRADENO\OUTREFUNDNO\REFUNDID
//orderIdType:订单ID类型，取值为TRANSACTIONID\OUTTRADENO\OUTREFUNDNO\REFUNDID之一
func (w WxPay) RefundQuery(orderId string, orderIdType ORDERIDTYPE) (map[string]string, error) {
	return newPayOrder().refundQuery(orderId, orderIdType)
}

//下载对账单
//bill_date：下载对账单的日期
//bill_type: 对账类型
func (w WxPay) DownLoadBill(bill_date string, bill_type BILLTYPE) (orderMap []map[string]string, sumMap map[string]string, err error) {
	return newPayOrder().download(bill_date, bill_type)
}

//支付结果通知
func (w WxPay) PayNotify(resp http.ResponseWriter, req *http.Request) {
	newPayOrder().payNotify(resp, req)
}

//退款通知
func (w WxPay) RefundNotify(resp http.ResponseWriter, req *http.Request) {
	newPayOrder().refundNotify(resp, req)
}

func logs(logType string, logContent string) {
	if logCallBack != nil {
		logCallBack.OnLog(logType, logContent)
	}
}
