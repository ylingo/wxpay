package wxpay

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func (o *PayOrder) download(bill_date string, bill_type BILLTYPE) (orderMap []map[string]string, sumMap map[string]string, err error) {
	reqEntry := download_Req{
		AppId:    cfg.App_Id,
		MchId:    cfg.Mch_Id,
		NonceStr: randstr.GetRandString(),
		SignType: string(cfg.DefaultSignType),
		BillDate: bill_date,
		BillType: string(bill_type),
		TarType:  "",
	}
	signXml := newSign().sign(reqEntry, cfg.Key, cfg.DefaultSignType)
	//fmt.Print(signXml)
	webReq := newHttpRequest()
	req, _ := webReq.getRequest("POST", cfg.WXBaseUrl+cfg.DownloadUrl, signXml)
	req.Header.Add("ContentType", "text/xml")
	var response *http.Response
	if response, err = webReq.do(req); err == nil {
		defer response.Body.Close()
		var bodyXml []byte
		if bodyXml, err = ioutil.ReadAll(response.Body); err == nil {
			if strings.Index(string(bodyXml), "<xml>") >= 0 &&
				strings.Index(string(bodyXml), "</xml>") > 0 {
				resp := make(map[string]interface{})
				if err = xml.Unmarshal(bodyXml, &resp); err == nil {
					return nil, nil, fmt.Errorf("%s,%s", resp["err_code"].(string), resp["err_code_des"].(string))
				}
			} else {
				records := [][]string{}
				reader := csv.NewReader(response.Body)
				for {
					record, err := reader.Read()
					if err == io.EOF {
						break
					}
					records = append(records, record)
				}
				recordCount := len(records)
				if recordCount > 3 {

					for i := 1; i < recordCount-2; i++ {
						m := make(map[string]string)
						fieldNum := len(records[0])
						for k := 0; k < fieldNum; k++ {
							m[records[0][k]] = records[i][k][1:]
						}
						orderMap = append(orderMap, m)
					}

					sumMap = make(map[string]string)
					fieldNum := len(records[recordCount-2])
					for k := 0; k < fieldNum; k++ {
						sumMap[records[recordCount-2][k]] = records[recordCount-1][k][1:]
					}
					return
				}
			}
		}
	}
	return
}
