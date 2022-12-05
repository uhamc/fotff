package gitee

import (
	"encoding/json"
	"fmt"
	"fotff/utils"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type PRIssueResp struct {
	URL string `json:"html_url"`
}

var mrIssueCache = cache.New(24*time.Hour, time.Hour)

func init() {
	if err := mrIssueCache.LoadFile("gitee_mr_issue.cache"); err != nil {
		fmt.Printf("load gitee_mr_issue.cache err: %v", err)
	}
}

func GetMRIssueURL(owner string, repo string, num int) (string, error) {
	url := fmt.Sprintf("https://gitee.com/api/v5/repos/%s/%s/pulls/%d/issues", owner, repo, num)
	var resp []byte
	if c, found := respCache.Get(url); found {
		resp = c.([]byte)
	} else {
		var err error
		resp, err = utils.DoSimpleHttpReq(http.MethodGet, url, nil)
		if err != nil {
			return "", err
		}
		respCache.Add(url, resp, cache.DefaultExpiration)
		respCache.SaveFile("gitee.cache")
	}
	var prIssue []PRIssueResp
	if err := json.Unmarshal(resp, &prIssue); err != nil {
		return "", err
	}
	if len(prIssue) == 0 {
		return "", nil
	}
	if len(prIssue) > 1 {
		logrus.Warnf("warn: find more than one issue related to %s, use the first one", url)
	}
	return prIssue[0].URL, nil
}
