package gitee

type CompareParam struct {
	Head  string
	Base  string
	Repo  string
	Owner string
}

type CommitDetail struct {
	Commit Commit `json:"commit"`
}

type Commit struct {
	Committer Committer `json:"committer"`
}

type Committer struct {
	Date string `json:"date"`
}

func GetBetweenMRs(param CompareParam) ([]CommitDetail, error) {
	// see: https://gitee.com/api/v5/swagger#/getV5ReposOwnerRepoCompareBase...Head
	panic("implement me")
}
