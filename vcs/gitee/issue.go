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
	mrIssueCache.LoadFile("gitee_mr_issue.cache")
}

func GetMRIssueURL(owner string, repo string, num int) (string, error) {
	defer mrIssueCache.SaveFile("gitee_mr_issue.cache")
	url := fmt.Sprintf("https://gitee.com/api/v5/repos/%s/%s/pulls/%d/issues", owner, repo, num)
	if c, found := mrIssueCache.Get(url); found {
		return c.(string), nil
	}
	resp, err := utils.DoSimpleHttpReq(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	var prIssue []PRIssueResp
	if err := json.Unmarshal(resp, &prIssue); err != nil {
		return "", err
	}
	if len(prIssue) == 0 {
		mrIssueCache.Add(url, "", cache.DefaultExpiration)
		return "", nil
	}
	if len(prIssue) > 1 {
		logrus.Warnf("warn: find more than one issue related to %s, use the first one", url)
	}
	mrIssueCache.Add(url, prIssue[0].URL, cache.DefaultExpiration)
	return prIssue[0].URL, nil
}
