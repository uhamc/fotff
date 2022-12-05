package gitee

import (
	"encoding/json"
	"fmt"
	"fotff/utils"
	"log"
	"net/http"
)

type PRIssueResp struct {
	URL string `json:"url"`
}

func GetMRIssueURL(owner string, repo string, num int) (string, error) {
	url := fmt.Sprintf("https://gitee.com/api/v5/repos/%s/%s/pulls/%d/issues", owner, repo, num)
	resp, err := utils.DoSimpleHttpReq(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	var prIssue []PRIssueResp
	if err := json.Unmarshal(resp, &prIssue); err != nil {
		return "", err
	}
	if len(prIssue) == 0 {
		return "", nil
	}
	if len(prIssue) > 1 {
		log.Printf("warn: find more than one issue related to %s, use the first one", url)
	}
	return prIssue[0].URL, nil
}
