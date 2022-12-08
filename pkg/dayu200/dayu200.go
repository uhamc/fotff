package dayu200

import (
	"code.cloudfoundry.org/archiver/extractor"
	"fmt"
	"fotff/pkg"
	"fotff/utils"
	"fotff/vcs"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type BuildServerConfig struct {
	Addr   string `key:"build_server_addr" default:"127.0.0.1:22"`
	User   string `key:"build_server_user" default:"root"`
	Passwd string `key:"build_server_password" default:"root"`
	// BuildWorkSpace must be absolute
	BuildWorkSpace string `key:"build_server_workspace" default:"/root/fotff/build_workspace"`
}

type Manager struct {
	ArchiveDir        string `key:"archive_dir" default:"."`
	Workspace         string `key:"workspace" default:"."`
	Branch            string `key:"branch" default:"master"`
	BuildServerConfig BuildServerConfig
	FlashTool         string `key:"flash_tool" default:"python"`
	FlashToolArgs     string `key:"flash_tool_args" default:"./flash.py"`
	lastFile          string
}

var stepCache = cache.New(24*time.Hour, time.Hour)

func init() {
	stepCache.LoadFile("dayu200_steps.cache")
}

func NewManager() pkg.Manager {
	var ret Manager
	utils.ParseFromConfigFile("dayu200", &ret)
	return &ret
}

func (m *Manager) Flash(pkg string) error {
	logrus.Infof("now flash %s", pkg)
	if _, err := os.Stat(filepath.Join(pkg, "__built__")); err != nil {
		if err := m.build(pkg); err != nil {
			logrus.Errorf("build pkg %s err: %v", pkg, err)
			return err
		}
	}
	logrus.Infof("calling flash tool for %s...", pkg)
	cmd := exec.Command(m.FlashTool, append(strings.Fields(m.FlashToolArgs), pkg)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Errorf("%s", string(out))
		logrus.Errorf("flash device fail: %v", err)
		return err
	}
	logrus.Infof("flash device successfully")
	return nil
}

func (m *Manager) Steps(from, to string) (pkgs []string, err error) {
	if from == to {
		return nil, fmt.Errorf("steps err: 'from' %s and 'to' %s are the same", from, to)
	}
	if c, found := stepCache.Get(from + "__to__" + to); found {
		logrus.Infof("steps from %s to %s are cached", from, to)
		return c.([]string), nil
	}
	defer stepCache.SaveFile("dayu200_steps.cache")
	updates, err := getRepoUpdates(from, to)
	if err != nil {
		return nil, err
	}
	logrus.Infof("find %d repo updates from %s to %s", len(updates), from, to)
	startTime, err := getPackageTime(from)
	if err != nil {
		return nil, err
	}
	endTime, err := getPackageTime(to)
	if err != nil {
		return nil, err
	}
	steps, err := getAllSteps(startTime, endTime, m.Branch, updates)
	if err != nil {
		return nil, err
	}
	logrus.Infof("find total %d steps from %s to %s", len(steps), from, to)
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
	stepCache.SaveFile("dayu200_steps.cache")
	return pkgs, nil
}

func (m *Manager) LastIssue(pkg string) (string, error) {
	//TODO implement me
	data, err := os.ReadFile(filepath.Join(pkg, "__last_issue__"))
	return string(data), err
}

func (m *Manager) GetNewer() (string, error) {
	m.lastFile = pkg.GetNewerFileFromDir(m.ArchiveDir, m.lastFile)
	ex := extractor.NewTgz()
	dirName := m.lastFile
	for filepath.Ext(dirName) != "" {
		dirName = strings.TrimSuffix(dirName, filepath.Ext(dirName))
	}
	dir := filepath.Join(m.Workspace, dirName)
	if _, err := os.Stat(dir); err == nil {
		return dir, nil
	}
	logrus.Infof("extracting %s to %s...", filepath.Join(m.ArchiveDir, m.lastFile), dir)
	if err := ex.Extract(filepath.Join(m.ArchiveDir, m.lastFile), dir); err != nil {
		return dir, err
	}
	if err := os.WriteFile(filepath.Join(dir, "__built__"), nil, 0640); err != nil {
		return "", err
	}
	return dir, nil
}
