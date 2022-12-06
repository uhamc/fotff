package dayu200

import (
	"fmt"
	"fotff/vcs"
	"fotff/vcs/gitee"
	"github.com/huandu/go-clone"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Step struct {
	IssueURLs []string
	MRs       []*gitee.Commit
}

func getRepoUpdates(from, to string) (updates []vcs.ProjectUpdate, err error) {
	m1, err := vcs.ParseManifestFile(filepath.Join(from, "manifest_tag.xml"))
	if err != nil {
		return nil, err
	}
	m2, err := vcs.ParseManifestFile(filepath.Join(to, "manifest_tag.xml"))
	if err != nil {
		return nil, err
	}
	return vcs.GetRepoUpdates(m1, m2, func(p1, p2 *vcs.Project) time.Time {
		logrus.Errorf("manifest structure changes not supported yet")
		return time.Time{}
	})
}

func getAllSteps(updates []vcs.ProjectUpdate) (ret []Step, err error) {
	allMRs, err := getAllMRs(updates)
	if err != nil {
		return nil, err
	}
	logrus.Infof("find total %d merge request commits of all repo updates", len(allMRs))
	issueInfos, err := combineMRsToIssue(allMRs)
	if err != nil {
		return nil, err
	}
	logrus.Infof("find total %d issues of all repo updates", len(issueInfos))
	return combineIssuesToStep(issueInfos)
}

func getAllMRs(updates []vcs.ProjectUpdate) (allMRs []*gitee.Commit, err error) {
	for _, update := range updates {
		//TODO remove this restrict
		if update.P1 == nil || update.P2 == nil {
			return nil, fmt.Errorf("find some repos added or removed, manifest structure changes not supported yet")
		}
		prs, err := gitee.GetBetweenMRs(gitee.CompareParam{
			Head:  update.P2.Revision,
			Base:  update.P1.Revision,
			Owner: "openharmony",
			Repo:  update.P2.Name,
		})
		if err != nil {
			return nil, err
		}
		allMRs = append(allMRs, prs...)
	}
	return
}

type IssueInfo struct {
	Visited       bool
	MRs           []*gitee.Commit
	RelatedIssues []string
}

func combineMRsToIssue(allMRs []*gitee.Commit) (map[string]*IssueInfo, error) {
	ret := make(map[string]*IssueInfo)
	for _, mr := range allMRs {
		num, err := strconv.Atoi(strings.Trim(regexp.MustCompile(`!\d+ `).FindString(mr.Commit.Message), "! "))
		if err != nil {
			return nil, fmt.Errorf("parse MR message for %s fail: %s", mr.URL, err)
		}
		issues, err := gitee.GetMRIssueURL(mr.Owner, mr.Repo, num)
		if err != nil {
			return nil, err
		}
		if len(issues) == 0 {
			issues = []string{mr.URL}
		}
		for i, issue := range issues {
			if _, ok := ret[issue]; !ok {
				ret[issue] = &IssueInfo{
					MRs:           []*gitee.Commit{mr},
					RelatedIssues: append(issues[:i], issues[i+1:]...),
				}
			} else {
				ret[issue] = &IssueInfo{
					MRs:           append(ret[issue].MRs, mr),
					RelatedIssues: append(ret[issue].RelatedIssues, append(issues[:i], issues[i+1:]...)...),
				}
			}
		}
	}
	return ret, nil
}

func combineOtherRelatedIssue(info *IssueInfo, all map[string]*IssueInfo) (mrs []*gitee.Commit, issues []string) {
	if info.Visited {
		return nil, nil
	}
	mrs = info.MRs
	issues = info.RelatedIssues
	info.Visited = true
	for _, other := range info.RelatedIssues {
		if i, ok := all[other]; ok {
			otherMRs, otherIssues := combineOtherRelatedIssue(i, all)
			mrs = append(mrs, otherMRs...)
			issues = append(issues, otherIssues...)
		}
		delete(all, other)
	}
	sort.Slice(mrs, func(i, j int) bool {
		// move the latest MR to the first place, use its merged_time to represent the update time of the issue
		return mrs[i].Commit.Committer.Date > mrs[j].Commit.Committer.Date
	})
	return deDupMRs(mrs), deDupIssues(issues)
}

func deDupMRs(mrs []*gitee.Commit) (retMRs []*gitee.Commit) {
	tmp := make(map[string]*gitee.Commit)
	for _, m := range mrs {
		tmp[m.SHA] = m
	}
	for _, m := range tmp {
		retMRs = append(retMRs, m)
	}
	return
}

func deDupIssues(issues []string) (retIssues []string) {
	tmp := make(map[string]string)
	for _, i := range issues {
		tmp[i] = i
	}
	for _, i := range tmp {
		retIssues = append(retIssues, i)
	}
	return
}

func combineIssuesToStep(issueInfos map[string]*IssueInfo) (ret []Step, err error) {
	for _, info := range issueInfos {
		info.MRs, info.RelatedIssues = combineOtherRelatedIssue(info, issueInfos)
	}
	for issue, infos := range issueInfos {
		ret = append(ret, Step{IssueURLs: append(infos.RelatedIssues, issue), MRs: infos.MRs})
	}
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].MRs[0].Commit.Committer.Date < ret[j].MRs[0].Commit.Committer.Date
	})
	return
}

func (m *Manager) genStepPackage(base *vcs.Manifest, step Step) (newPkg string, newManifest *vcs.Manifest, err error) {
	defer func() {
		logrus.Infof("package dir %s for step %v generated", newPkg, step.IssueURLs)
	}()
	newManifest = clone.Clone(base).(*vcs.Manifest)
	for _, mr := range step.MRs {
		newManifest.UpdateManifestProject(mr.Repo, "", "", mr.SHA)
	}
	md5sum, err := newManifest.Standardize()
	if err != nil {
		return "", nil, err
	}
	if _, err := os.Stat(filepath.Join(m.Workspace, md5sum, "__built__")); err == nil {
		return filepath.Join(m.Workspace, md5sum), newManifest, nil
	}
	if err := os.MkdirAll(filepath.Join(m.Workspace, md5sum), 0750); err != nil {
		return "", nil, err
	}
	if err := os.WriteFile(filepath.Join(m.Workspace, md5sum, "__last_issue__"), []byte(fmt.Sprintf("%v", step.IssueURLs)), 0640); err != nil {
		return "", nil, err
	}
	err = newManifest.WriteFile(filepath.Join(m.Workspace, md5sum, "manifest_tag.xml"))
	if err != nil {
		return "", nil, err
	}
	return filepath.Join(m.Workspace, md5sum), newManifest, nil
}
