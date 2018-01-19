微信支付API封装

微信支付统一下单、查询订单、关闭订单、退费、退费查询、下载订单，以及支付通知、退费通知等接口
所有请求接口返回map[string]interface{}
需要实现通知逻辑接口函数
        OnPayNotify(map[string]interface{})
	OnRefundNotify(map[string]interface{})
当微信通知到达后会调用该接口