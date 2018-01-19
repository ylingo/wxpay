package wxpay

type LogCallBack interface {
	OnLog(logType string, logContent string)
}

type NotifyCallBack interface {
	OnNotify()
}

type NotifyLogic interface {
	OnPayNotify(map[string]interface{})
	OnRefundNotify(map[string]interface{})
}
