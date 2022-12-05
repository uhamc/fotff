package gitee

import (
	"encoding/json"
	"fmt"
	"fotff/utils"
	"github.com/patrickmn/go-cache"
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
	Commits []Commit `json:"commits"`
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

var compareCache = cache.New(24*time.Hour, time.Hour)

func init() {
	compareCache.LoadFile("gitee_compare.cache")
}

func GetBetweenMRs(param CompareParam) ([]Commit, error) {
	commits, err := GetBetweenCommits(param)
	if err != nil {
		return nil, err
	}
	var ret []Commit
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

func GetBetweenCommits(param CompareParam) ([]Commit, error) {
	defer compareCache.SaveFile("gitee_compare.cache")
	url := fmt.Sprintf("https://gitee.com/api/v5/repos/%s/%s/compare/%s...%s", param.Owner, param.Repo, param.Base, param.Head)
	if c, found := compareCache.Get(url); found {
		return c.([]Commit), nil
	}
	resp, err := utils.DoSimpleHttpReq(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	var compareResp CompareResp
	if err := json.Unmarshal(resp, &compareResp); err != nil {
		return nil, err
	}
	compareCache.Add(url, compareResp.Commits, cache.DefaultExpiration)
	return compareResp.Commits, nil
}
