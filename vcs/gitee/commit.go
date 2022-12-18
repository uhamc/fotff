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

package gitee

import (
	"encoding/json"
	"fmt"
	"fotff/utils"
	"net/http"
)

func GetCommit(owner, repo, id string) (*Commit, error) {
	url := fmt.Sprintf("https://gitee.com/api/v5/repos/%s/%s/commits/%s", owner, repo, id)
	var resp []byte
	if c, found := utils.CacheGet("gitee", url); found {
		resp = c.([]byte)
	} else {
		var err error
		resp, err = utils.DoSimpleHttpReq(http.MethodGet, url, nil, nil)
		if err != nil {
			return nil, err
		}
		utils.CacheSet("gitee", url, resp)
	}
	var commitResp Commit
	if err := json.Unmarshal(resp, &commitResp); err != nil {
		return nil, err
	}
	commitResp.Owner = owner
	commitResp.Repo = repo
	return &commitResp, nil
}
