package gitee

import (
	"encoding/json"
	"fmt"
	"fotff/utils"
	"net/http"
	"time"
)

type CompareParam struct {
	Head  string
	Base  string
	Repo  string
	Owner string
}

type CompareResp struct {
	Commits []*Commit `json:"commits"`
}

type Commit struct {
	CommitExtend `json:"-"`
	URL          string `json:"url"`
	SHA          string `json:"sha"`
	Commit       struct {
		Committer struct {
			Date string `json:"date"`
		} `json:"committer"`
		Message string `json:"message"`
	} `json:"commit"`
	Parents []struct {
		SHA string `json:"sha"`
		URL string `json:"url"`
	} `json:"parents"`
	Files []struct {
		Filename string `json:"filename"`
		Status   string `json:"status"`
		Patch    string `json:"patch,omitempty"`
	} `json:"files,omitempty"`
}

type CommitExtend struct {
	Owner string
	Repo  string
}

func GetLatestMRBefore(owner, repo, branch string, before string) (ret *Commit, err error) {
	branchResp, err := GetBranch(owner, repo, branch)
	if err != nil {
		return nil, err
	}
	head := branchResp.Commit
	head.Owner = owner
	head.Repo = repo
	for head.Commit.Committer.Date > before {
		if head, err = GetCommit(owner, repo, head.Parents[0].SHA); err != nil {
			return nil, err
		}
	}
	return head, nil
}

func GetBetweenTimeMRs(owner, repo, branch string, from, to time.Time) (ret []*Commit, err error) {
	branchResp, err := GetBranch(owner, repo, branch)
	if err != nil {
		return nil, err
	}
	fromStr := from.UTC().Format(time.RFC3339)
	toStr := to.UTC().Format(time.RFC3339)
	head := branchResp.Commit
	head.Owner = owner
	head.Repo = repo
	for head.Commit.Committer.Date > fromStr {
		if head.Commit.Committer.Date < toStr {
			ret = append(ret, head)
		}
		if head, err = GetCommit(owner, repo, head.Parents[0].SHA); err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func GetBetweenMRs(param CompareParam) ([]*Commit, error) {
	commits, err := GetBetweenCommits(param)
	if err != nil {
		return nil, err
	}
	var ret []*Commit
	head := param.Head
	for head != param.Base {
		for _, commit := range commits {
			if commit.SHA != head {
				continue
			}
			commit.Owner = param.Owner
			commit.Repo = param.Repo
			ret = append(ret, commit)
			head = commit.Parents[0].SHA
		}
	}
	return ret, nil
}

func GetBetweenCommits(param CompareParam) ([]*Commit, error) {
	url := fmt.Sprintf("https://gitee.com/api/v5/repos/%s/%s/compare/%s...%s", param.Owner, param.Repo, param.Base, param.Head)
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
	var compareResp CompareResp
	if err := json.Unmarshal(resp, &compareResp); err != nil {
		return nil, err
	}
	return compareResp.Commits, nil
}
