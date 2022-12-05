package dayu200

import (
	"code.cloudfoundry.org/archiver/extractor"
	"fmt"
	"fotff/pkg"
	"fotff/vcs"
	"github.com/patrickmn/go-cache"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Manager struct {
	PkgDir    string
	Workspace string
	lastFile  string
}

var stepCache = cache.New(24*time.Hour, time.Hour)

func (m *Manager) Flash(pkg string) error {
	//TODO implement me
	panic("implement me")
}

func (m *Manager) Steps(from, to string) (pkgs []string, err error) {
	if from == to {
		return nil, fmt.Errorf("steps err: 'from' %s and 'to' %s are the same", from, to)
	}
	if c, found := stepCache.Get(from + "__to__" + to); found {
		return c.([]string), nil
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
	stepCache.Add(from+"__to__"+to, pkgs, cache.DefaultExpiration)
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
	dirName := m.lastFile
	for filepath.Ext(dirName) != "" {
		dirName = strings.TrimSuffix(dirName, filepath.Ext(dirName))
	}
	dir := filepath.Join(m.Workspace, dirName)
	if _, err := os.Stat(dir); err == nil {
		return dir, nil
	}
	log.Printf("extracting %s to %s...", filepath.Join(m.PkgDir, m.lastFile), dir)
	if err := ex.Extract(filepath.Join(m.PkgDir, m.lastFile), dir); err != nil {
		return dir, err
	}
	if err := os.WriteFile(filepath.Join(dir, "__built__"), nil, 0640); err != nil {
		return "", err
	}
	return dir, nil
}
