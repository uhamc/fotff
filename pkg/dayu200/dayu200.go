package dayu200

import (
	"code.cloudfoundry.org/archiver/extractor"
	"fotff/pkg"
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
	//TODO implement me
	panic("implement me")
}

func (m *Manager) LastIssue(pkg string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *Manager) GetNewer() (string, error) {
	m.lastFile = pkg.GetNewerFileFromDir(m.PkgDir, m.lastFile)
	ex := extractor.NewTgz()
	err := ex.Extract(filepath.Join(m.PkgDir, m.lastFile), filepath.Join(m.Workspace, strings.TrimSuffix(m.lastFile, filepath.Ext(m.lastFile))))
	return m.Workspace, err
}
