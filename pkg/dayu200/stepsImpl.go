package dayu200

import (
	"fmt"
	"fotff/vcs"
	"fotff/vcs/gitee"
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
	IssueURL string
	MRs      []gitee.Commit
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
	issueMRs, err := combineMRsToIssue(allMRs)
	if err != nil {
		return nil, err
	}
	logrus.Infof("find total %d issues of all repo updates, use each issue as one step", len(issueMRs))
	for issue, mrs := range issueMRs {
		ret = append(ret, Step{IssueURL: issue, MRs: mrs})
	}
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].MRs[0].Commit.Committer.Date > ret[j].MRs[0].Commit.Committer.Date
	})
	return ret, err
}

func getAllMRs(updates []vcs.ProjectUpdate) (allMRs []gitee.Commit, err error) {
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

func combineMRsToIssue(allMRs []gitee.Commit) (map[string][]gitee.Commit, error) {
	ret := make(map[string][]gitee.Commit)
	for _, mr := range allMRs {
		num, err := strconv.Atoi(strings.Trim(regexp.MustCompile(`!\d+ `).FindString(mr.Commit.Message), "! "))
		if err != nil {
			return nil, fmt.Errorf("parse MR message for %s fail: %s", mr.URL, err)
		}
		issue, err := gitee.GetMRIssueURL(mr.Owner, mr.Repo, num)
		if err != nil {
			return nil, err
		}
		if issue == "" {
			issue = mr.URL
		}
		mrList := append(ret[issue], mr)
		sort.Slice(mrList, func(i, j int) bool {
			// move the latest MR to the first place, use its merged_time to represent the update time of the issue
			return mrList[i].Commit.Committer.Date > mrList[j].Commit.Committer.Date
		})
		ret[issue] = mrList
	}
	return ret, nil
}

func (m *Manager) genStepPackage(base *vcs.Manifest, step Step) (newPkg string, newManifest *vcs.Manifest, err error) {
	defer func() {
		logrus.Infof("package dir %s for step %s generated", newPkg, step.IssueURL)
	}()
	newManifest = base.DeepCopy()
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
	if err := os.WriteFile(filepath.Join(m.Workspace, md5sum, "__last_issue__"), []byte(step.IssueURL), 0640); err != nil {
		return "", nil, err
	}
	err = newManifest.WriteFile(filepath.Join(m.Workspace, md5sum, "manifest_tag.xml"))
	if err != nil {
		return "", nil, err
	}
	return filepath.Join(m.Workspace, md5sum), newManifest, nil
}
