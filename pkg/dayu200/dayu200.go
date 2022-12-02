package dayu200

import (
	"code.cloudfoundry.org/archiver/extractor"
	"fmt"
	"fotff/pkg"
	"fotff/vcs"
	"fotff/vcs/gitee"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var simpleRegTimeInPkgName = regexp.MustCompile(`\d{8}_\d{6}`)

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
	changes, adds, removes, err := vcs.GetChangedRepos(m1, m2)
	if err != nil {
		return nil, err
	}
	if len(changes)+len(adds)+len(removes) == 0 {
		return []string{from, to}, nil
	}
	t1, err := getPackageTime(from)
	if err != nil {
		return nil, err
	}
	t2, err := getPackageTime(to)
	if err != nil {
		return nil, err
	}
	panic("implement me")
	for range changes {
		gitee.GetPRs(gitee.PRSearchParam{
			Since: t1,
			ExtendPRSearchParam: gitee.ExtendPRSearchParam{
				Until: t2,
			},
		})
	}
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

func getPackageTime(pkg string) (time.Time, error) {
	return time.ParseInLocation(`20060102_150405`, simpleRegTimeInPkgName.FindString(filepath.Base(pkg)), time.Local)
}
