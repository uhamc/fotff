package utils

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

func DoSimpleHttpReq(method string, url string, body []byte) (ret []byte, err error) {
	for i := 0; i < 3; i++ {
		if ret, err = doSimpleHttpReqImpl(method, url, body); err == nil {
			return
		}
		time.Sleep(time.Second)
	}
	return
}

func doSimpleHttpReqImpl(method string, url string, body []byte) (ret []byte, err error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	resp, err := proxyClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		data, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == http.StatusProxyAuthRequired || resp.StatusCode == http.StatusForbidden {
			SwitchProxy()
		}
		logrus.Errorf("%s %s: code: %d body: %s", method, url, resp.StatusCode, string(data))
		return nil, fmt.Errorf("%s %s: code: %d body: %s", method, url, resp.StatusCode, string(data))
	}
	return io.ReadAll(resp.Body)
}
