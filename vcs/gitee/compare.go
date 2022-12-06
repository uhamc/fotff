package gitee

import (
	"encoding/json"
	"fmt"
	"fotff/utils"
	"github.com/patrickmn/go-cache"
	"net/http"
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
}

type CommitExtend struct {
	Owner string
	Repo  string
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
	var compareResp CompareResp
	if err := json.Unmarshal(resp, &compareResp); err != nil {
		return nil, err
	}
	return compareResp.Commits, nil
}
