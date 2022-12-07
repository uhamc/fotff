package gitee

import (
	"encoding/json"
	"fmt"
	"fotff/utils"
	"github.com/patrickmn/go-cache"
	"net/http"
)

func GetCommit(owner, repo, id string) (*Commit, error) {
	url := fmt.Sprintf("https://gitee.com/api/v5/repos/%s/%s/commits/%s", owner, repo, id)
	var resp []byte
	if c, found := respCache.Get(url); found {
		resp = c.([]byte)
	} else {
		var err error
		resp, err = utils.DoSimpleHttpReq(http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
		respCache.Add(url, resp, cache.DefaultExpiration)
		respCache.SaveFile("gitee.cache")
	}
	var commitResp Commit
	if err := json.Unmarshal(resp, &commitResp); err != nil {
		return nil, err
	}
	return &commitResp, nil
}
