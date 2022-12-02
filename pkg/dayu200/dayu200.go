package dayu200

import (
	"code.cloudfoundry.org/archiver/extractor"
	"fmt"
	"fotff/pkg"
	"fotff/vcs"
	"fotff/vcs/gitee"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Manager struct {
	PkgDir    string
	Workspace string
	lastFile  string
}

func (m *Manager) Flash(pkg string) error {
	//TODO implement me
	panic("implement me")
}

func (m *Manager) Steps(from, to string) ([]string, error) {
	if from == to {
		return nil, fmt.Errorf("steps err: 'from' %s and 'to' %s are the same", from, to)
	}
	m1, err := vcs.ParseManifestFile(filepath.Join(from, "manifest_tag.xml"))
	if err != nil {
		return nil, err
	}
	m2, err := vcs.ParseManifestFile(filepath.Join(from, "manifest_tag.xml"))
	if err != nil {
		return nil, err
	}
	updates, err := vcs.GetChangedRepos(m1, m2)
	if err != nil {
		return nil, err
	}
	if len(updates) == 0 {
		return []string{from, to}, nil
	}
	var allMRs []gitee.CommitDetail
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
	sort.SliceStable(allMRs, func(i, j int) bool {
		return allMRs[i].Commit.Committer.Date < allMRs[j].Commit.Committer.Date
	})
	//TODO generate manifests
	panic("implement me")
}

func (m *Manager) LastIssue(pkg string) (string, error) {
	//TODO implement me
	data, err := os.ReadFile(filepath.Join(pkg, "__last_issue__"))
	return string(data), err
}

func (m *Manager) GetNewer() (string, error) {
	m.lastFile = pkg.GetNewerFileFromDir(m.PkgDir, m.lastFile)
	ex := extractor.NewTgz()
	err := ex.Extract(filepath.Join(m.PkgDir, m.lastFile), filepath.Join(m.Workspace, strings.TrimSuffix(m.lastFile, filepath.Ext(m.lastFile))))
	return m.Workspace, err
}

func build(pkg string) error {
	if _, err := os.Stat(filepath.Join(pkg, "__to_be_built__")); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	//TODO build package with generated manifest_tag.xml
	panic("implement me")
}
