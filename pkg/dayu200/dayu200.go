package dayu200

import (
	"code.cloudfoundry.org/archiver/extractor"
	"fmt"
	"fotff/pkg"
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

func (m *Manager) Steps(from, to string) ([]string, error) {
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
	//TODO generate manifests
	panic("implement me")
	panic(steps)
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
