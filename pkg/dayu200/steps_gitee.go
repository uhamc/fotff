package dayu200

import (
	"bufio"
	"bytes"
	"encoding/xml"
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
	"sync"
	"time"
)

type IssueInfo struct {
	visited          bool
	RelatedIssues    []string
	MRs              []*gitee.Commit
	StructCTime      string
	StructureUpdates []*vcs.ProjectUpdate
}

type Step struct {
	IssueURLs        []string
	MRs              []*gitee.Commit
	StructCTime      string
	StructureUpdates []*vcs.ProjectUpdate
}

func (m *Manager) stepsFromGitee(from, to string) (pkgs []string, err error) {
	updates, err := m.getRepoUpdates(from, to)
	if err != nil {
		return nil, err
	}
	startTime, err := getPackageTime(from)
	if err != nil {
		return nil, err
	}
	endTime, err := getPackageTime(to)
	if err != nil {
		return nil, err
	}
	logrus.Infof("find %d repo updates from %s to %s", len(updates), from, to)
	steps, err := getAllStepsFromGitee(startTime, endTime, m.Branch, updates)
	if err != nil {
		return nil, err
	}
	logrus.Infof("find total %d steps from %s to %s", len(steps), from, to)
	baseManifest, err := vcs.ParseManifestFile(filepath.Join(m.Workspace, from, "manifest_tag.xml"))
	if err != nil {
		return nil, err
	}
	for _, step := range steps {
		var newPkg string
		if newPkg, baseManifest, err = m.genStepPackage(baseManifest, step); err != nil {
			return nil, err
		}
		pkgs = append(pkgs, newPkg)
	}
	return pkgs, nil
}

func (m *Manager) getRepoUpdates(from, to string) (updates []vcs.ProjectUpdate, err error) {
	m1, err := vcs.ParseManifestFile(filepath.Join(m.Workspace, from, "manifest_tag.xml"))
	if err != nil {
		return nil, err
	}
	m2, err := vcs.ParseManifestFile(filepath.Join(m.Workspace, to, "manifest_tag.xml"))
	if err != nil {
		return nil, err
	}
	return vcs.GetRepoUpdates(m1, m2)
}

func getAllStepsFromGitee(startTime, endTime time.Time, branch string, updates []vcs.ProjectUpdate) (ret []Step, err error) {
	allMRs, err := getAllMRs(startTime, endTime, branch, updates)
	if err != nil {
		return nil, err
	}
	issueInfos, err := combineMRsToIssue(allMRs, branch)
	if err != nil {
		return nil, err
	}
	return combineIssuesToStep(issueInfos)
}

func getAllMRs(startTime, endTime time.Time, branch string, updates []vcs.ProjectUpdate) (allMRs []*gitee.Commit, err error) {
	var once sync.Once
	for _, update := range updates {
		var prs []*gitee.Commit
		if update.P1.StructureDiff(update.P2) {
			once.Do(func() {
				prs, err = gitee.GetBetweenTimeMRs("openharmony", "manifest", branch, startTime, endTime)
			})
			if update.P1 != nil {
				var p1 []*gitee.Commit
				p1, err = gitee.GetBetweenTimeMRs("openharmony", update.P1.Name, branch, startTime, endTime)
				prs = append(prs, p1...)
			}
			if update.P2 != nil {
				var p2 []*gitee.Commit
				p2, err = gitee.GetBetweenTimeMRs("openharmony", update.P2.Name, branch, startTime, endTime)
				prs = append(prs, p2...)
			}
		} else {
			prs, err = gitee.GetBetweenMRs(gitee.CompareParam{
				Head:  update.P2.Revision,
				Base:  update.P1.Revision,
				Owner: "openharmony",
				Repo:  update.P2.Name,
			})
		}
		if err != nil {
			return nil, err
		}
		allMRs = append(allMRs, prs...)
	}
	logrus.Infof("find total %d merge request commits of all repo updates", len(allMRs))
	return
}

func combineMRsToIssue(allMRs []*gitee.Commit, branch string) (map[string]*IssueInfo, error) {
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
		var scs []*vcs.ProjectUpdate
		var scTime string
		if mr.Owner == "openharmony" && mr.Repo == "manifest" {
			if scTime, scs, err = parseStructureUpdates(mr, branch); err != nil {
				return nil, err
			}
		}
		for i, issue := range issues {
			if _, ok := ret[issue]; !ok {
				ret[issue] = &IssueInfo{
					MRs:              []*gitee.Commit{mr},
					RelatedIssues:    append(issues[:i], issues[i+1:]...),
					StructCTime:      scTime,
					StructureUpdates: scs,
				}
			} else {
				ret[issue] = &IssueInfo{
					MRs:              append(ret[issue].MRs, mr),
					RelatedIssues:    append(ret[issue].RelatedIssues, append(issues[:i], issues[i+1:]...)...),
					StructCTime:      scTime,
					StructureUpdates: append(ret[issue].StructureUpdates, scs...),
				}
			}
		}
	}
	logrus.Infof("find total %d issues of all repo updates", len(ret))
	return ret, nil
}

func combineOtherRelatedIssue(parent, self *IssueInfo, all map[string]*IssueInfo) {
	if self.visited {
		return
	}
	self.visited = true
	for _, other := range self.RelatedIssues {
		if son, ok := all[other]; ok {
			combineOtherRelatedIssue(self, son, all)
			delete(all, other)
		}
	}
	parent.RelatedIssues = deDupIssues(append(parent.RelatedIssues, self.RelatedIssues...))
	parent.MRs = deDupMRs(append(parent.MRs, self.MRs...))
	parent.StructureUpdates = deDupProjectUpdates(append(parent.StructureUpdates, self.StructureUpdates...))
	if len(parent.StructCTime) != 0 && parent.StructCTime < self.StructCTime {
		parent.StructCTime = self.StructCTime
	}
}

func deDupProjectUpdates(us []*vcs.ProjectUpdate) (retMRs []*vcs.ProjectUpdate) {
	dupIndexes := make([]bool, len(us))
	for i := range us {
		for j := i + 1; j < len(us); j++ {
			if us[j].P1 == us[i].P1 && us[j].P2 == us[i].P2 {
				dupIndexes[j] = true
			}
		}
	}
	for i, dup := range dupIndexes {
		if dup {
			continue
		}
		retMRs = append(retMRs, us[i])
	}
	return
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

// parseStructureUpdates get changed XMLs and parse it to recognize repo structure changes.
// Since we do not care which revision a repo was, P1 is not welly handled, just assign it not nil for performance.
func parseStructureUpdates(commit *gitee.Commit, branch string) (string, []*vcs.ProjectUpdate, error) {
	tmp := make(map[string]vcs.ProjectUpdate)
	if len(commit.Files) == 0 {
		// commit that queried from MR req does not contain file details, should fetch again
		var err error
		if commit, err = gitee.GetCommit(commit.Owner, commit.Repo, commit.SHA); err != nil {
			return "", nil, err
		}
	}
	for _, f := range commit.Files {
		if filepath.Ext(f.Filename) != ".xml" {
			continue
		}
		if err := parseFilePatch(f.Patch, tmp); err != nil {
			return "", nil, err
		}
	}
	var ret []*vcs.ProjectUpdate
	for _, pu := range tmp {
		projectUpdateCopy := pu
		ret = append(ret, &projectUpdateCopy)
	}
	for _, pu := range ret {
		if pu.P1 == nil && pu.P2 != nil {
			lastCommit, err := gitee.GetLatestMRBefore("openharmony", pu.P2.Name, branch, commit.Commit.Committer.Date)
			if err != nil {
				return "", nil, err
			}
			pu.P2.Revision = lastCommit.SHA
		}
	}
	return commit.Commit.Committer.Date, ret, nil
}

func parseFilePatch(str string, m map[string]vcs.ProjectUpdate) error {
	sc := bufio.NewScanner(bytes.NewBuffer([]byte(str)))
	for sc.Scan() {
		line := sc.Text()
		var p vcs.Project
		if strings.HasPrefix(line, "-") {
			if err := xml.Unmarshal([]byte(line[1:]), &p); err == nil {
				m[p.Name] = vcs.ProjectUpdate{P1: &p, P2: m[p.Name].P2}
			}
		} else if strings.HasPrefix(line, "+") {
			if err := xml.Unmarshal([]byte(line[1:]), &p); err == nil {
				m[p.Name] = vcs.ProjectUpdate{P1: m[p.Name].P1, P2: &p}
			}
		}
	}
	return nil
}

func combineIssuesToStep(issueInfos map[string]*IssueInfo) (ret []Step, err error) {
	for _, info := range issueInfos {
		combineOtherRelatedIssue(info, info, issueInfos)
	}
	for issue, infos := range issueInfos {
		sort.Slice(infos.MRs, func(i, j int) bool {
			// move the latest MR to the first place, use its merged_time to represent the update time of the issue
			return infos.MRs[i].Commit.Committer.Date > infos.MRs[j].Commit.Committer.Date
		})
		ret = append(ret, Step{
			IssueURLs:        append(infos.RelatedIssues, issue),
			MRs:              infos.MRs,
			StructCTime:      infos.StructCTime,
			StructureUpdates: infos.StructureUpdates})
	}
	sort.Slice(ret, func(i, j int) bool {
		ti, tj := ret[i].MRs[0].Commit.Committer.Date, ret[j].MRs[0].Commit.Committer.Date
		if len(ret[i].StructCTime) != 0 {
			ti = ret[i].StructCTime
		}
		if len(ret[j].StructCTime) != 0 {
			ti = ret[j].StructCTime
		}
		return ti < tj
	})
	logrus.Infof("find total %d steps of all issues", len(ret))
	return
}

var simpleRegTimeInPkgName = regexp.MustCompile(`\d{8}_\d{6}`)

func getPackageTime(pkg string) (time.Time, error) {
	return time.ParseInLocation(`20060102_150405`, simpleRegTimeInPkgName.FindString(pkg), time.Local)
}

func (m *Manager) genStepPackage(base *vcs.Manifest, step Step) (newPkg string, newManifest *vcs.Manifest, err error) {
	defer func() {
		logrus.Infof("package dir %s for step %v generated", newPkg, step.IssueURLs)
	}()
	newManifest = clone.Clone(base).(*vcs.Manifest)
	for _, u := range step.StructureUpdates {
		if u.P2 != nil {
			newManifest.UpdateManifestProject(u.P2.Name, u.P2.Path, u.P2.Remote, u.P2.Revision, true)
		} else if u.P1 != nil {
			newManifest.RemoveManifestProject(u.P1.Name)
		}
	}
	for _, mr := range step.MRs {
		newManifest.UpdateManifestProject(mr.Repo, "", "", mr.SHA, false)
	}
	md5sum, err := newManifest.Standardize()
	if err != nil {
		return "", nil, err
	}
	if _, err := os.Stat(filepath.Join(m.Workspace, md5sum, "__built__")); err == nil {
		return md5sum, newManifest, nil
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
	return md5sum, newManifest, nil
}
