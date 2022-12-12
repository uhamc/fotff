package gitee

import (
	"encoding/json"
	"fmt"
	"fotff/utils"
	"net/http"
)

type PRIssueResp struct {
	URL string `json:"html_url"`
}

func GetMRIssueURL(owner string, repo string, num int) ([]string, error) {
	url := fmt.Sprintf("https://gitee.com/api/v5/repos/%s/%s/pulls/%d/issues", owner, repo, num)
	var resp []byte
	if c, found := utils.CacheGet("gitee", url); found {
		resp = c.([]byte)
	} else {
		var err error
		resp, err = utils.DoSimpleHttpReq(http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
		utils.CacheSet("gitee", url, resp)
	}
	var prIssues []PRIssueResp
	if err := json.Unmarshal(resp, &prIssues); err != nil {
		return nil, err
	}
	ret := make([]string, len(prIssues))
	for i, issue := range prIssues {
		ret[i] = issue.URL
	}
	return ret, nil
}
