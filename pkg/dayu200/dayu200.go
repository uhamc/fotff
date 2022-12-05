package dayu200

import (
	"code.cloudfoundry.org/archiver/extractor"
	"fmt"
	"fotff/pkg"
	"fotff/vcs"
	"os"
	"path/filepath"
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

func (m *Manager) Steps(from, to string) (pkgs []string, err error) {
	if from == to {
		return nil, fmt.Errorf("steps err: 'from' %s and 'to' %s are the same", from, to)
	}
	updates, err := getRepoUpdates(from, to)
	if err != nil {
		return nil, err
	}
	steps, err := getAllSteps(updates)
	if err != nil {
		return nil, err
	}
	baseManifest, err := vcs.ParseManifestFile(filepath.Join(from, "manifest_tag.xml"))
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

func (m *Manager) LastIssue(pkg string) (string, error) {
	//TODO implement me
	data, err := os.ReadFile(filepath.Join(pkg, "__last_issue__"))
	return string(data), err
}

func (m *Manager) GetNewer() (string, error) {
	m.lastFile = pkg.GetNewerFileFromDir(m.PkgDir, m.lastFile)
	ex := extractor.NewTgz()
	dir := filepath.Join(m.Workspace, strings.TrimSuffix(m.lastFile, filepath.Ext(m.lastFile)))
	if err := os.RemoveAll(dir); err != nil {
		return dir, err
	}
	if err := ex.Extract(filepath.Join(m.PkgDir, m.lastFile), dir); err != nil {
		return dir, err
	}
	if err := os.WriteFile(filepath.Join(dir, "__built__"), nil, 0640); err != nil {
		return "", err
	}
	return dir, nil
}
