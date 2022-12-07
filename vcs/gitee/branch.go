package gitee

import (
	"encoding/json"
	"fmt"
	"fotff/utils"
	"net/http"
)

type BranchResp struct {
	Name   string  `json:"name"`
	Commit *Commit `json:"commit"`
}

func GetBranch(owner, repo, branch string) (*BranchResp, error) {
	url := fmt.Sprintf("https://gitee.com/api/v5/repos/%s/%s/branches/%s", owner, repo, branch)
	resp, err := utils.DoSimpleHttpReq(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	var branchResp BranchResp
	if err := json.Unmarshal(resp, &branchResp); err != nil {
		return nil, err
	}
	return &branchResp, nil
}
