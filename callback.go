package wxpay

type LogCallBack interface {
	OnLog(logType string, logContent string)
}

type NotifyLogic interface {
	OnPayNotify(map[string]string)
	OnRefundNotify(map[string]string)
}
