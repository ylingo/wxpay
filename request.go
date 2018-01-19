package wxpay

import (
	"net"
	"net/http"
	"strings"
	"time"
)

type httpRequest struct {
}

func newHttpRequest() *httpRequest {
	return &httpRequest{}
}

func (r httpRequest) getRequest(method, url, strBody string) (*http.Request, error) {
	return http.NewRequest(method, url, strings.NewReader(strBody))
}

func (r httpRequest) do(req *http.Request) (*http.Response, error) {
	return r.client().Do(req)
}

func (r httpRequest) client() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(25 * time.Second)
				c, err := net.DialTimeout(netw, addr, time.Second*20)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}
}
