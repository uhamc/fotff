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
