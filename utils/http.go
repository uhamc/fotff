/*
 * Copyright (c) 2022 Huawei Device Co., Ltd.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package utils

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

func DoSimpleHttpReqRaw(method string, url string, body []byte, header map[string]string) (response *http.Response, err error) {
	for i := 0; i < 3; i++ {
		if response, err = doSimpleHttpReqImpl(method, url, body, header); err == nil {
			return
		}
		time.Sleep(time.Second)
	}
	return
}

func DoSimpleHttpReq(method string, url string, body []byte, header map[string]string) (ret []byte, err error) {
	var resp *http.Response
	for i := 0; i < 3; i++ {
		if resp, err = doSimpleHttpReqImpl(method, url, body, header); err == nil {
			ret, err = io.ReadAll(resp.Body)
			resp.Body.Close()
			return
		}
		time.Sleep(time.Second)
	}
	return
}

func doSimpleHttpReqImpl(method string, url string, body []byte, header map[string]string) (response *http.Response, err error) {
	logrus.Infof("%s %s", method, url)
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	resp, err := proxyClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		defer resp.Body.Close()
		data, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == http.StatusProxyAuthRequired || resp.StatusCode == http.StatusForbidden {
			SwitchProxy()
		}
		logrus.Errorf("%s %s: code: %d body: %s", method, url, resp.StatusCode, string(data))
		return nil, fmt.Errorf("%s %s: code: %d body: %s", method, url, resp.StatusCode, string(data))
	}
	return resp, nil
}
