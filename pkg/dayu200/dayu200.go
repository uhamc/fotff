package dayu200

import (
	"code.cloudfoundry.org/archiver/extractor"
	"fmt"
	"fotff/pkg"
	"fotff/utils"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	SN                string `key:"sn" default:""`
	hdc               string
}

func NewManager() pkg.Manager {
	var ret Manager
	utils.ParseFromConfigFile("dayu200", &ret)
	if ret.hdc, _ = exec.LookPath("hdc"); ret.hdc == "" {
		ret.hdc, _ = exec.LookPath("hdc_std")
	}
	if ret.hdc == "" {
		logrus.Panicf("can not find 'hdc', please install")
	}
	return &ret
}

func (m *Manager) Flash(pkg string) error {
	logrus.Infof("now flash %s", pkg)
	if _, err := os.Stat(filepath.Join(m.Workspace, pkg, "__built__")); err != nil {
		if err := m.build(pkg); err != nil {
			logrus.Errorf("build pkg %s err: %v", pkg, err)
			return err
		}
	}
	return m.flashDevice(pkg)
}

func (m *Manager) Steps(from, to string) (pkgs []string, err error) {
	if from == to {
		return nil, fmt.Errorf("steps err: 'from' %s and 'to' %s are the same", from, to)
	}
	if c, found := utils.CacheGet("dayu200_steps", from+"__to__"+to); found {
		logrus.Infof("steps from %s to %s are cached", from, to)
		return c.([]string), nil
	}
	if pkgs, err = m.stepsFromGitee(from, to); err != nil {
		return pkgs, err
	}
	utils.CacheSet("dayu200_steps", from+"__to__"+to, pkgs)
	return pkgs, nil
}

func (m *Manager) LastIssue(pkg string) (string, error) {
	//TODO implement me
	data, err := os.ReadFile(filepath.Join(m.Workspace, pkg, "__last_issue__"))
	return string(data), err
}

func (m *Manager) GetNewer(cur string) (string, error) {
	newFile := pkg.GetNewerFileFromDir(m.ArchiveDir, cur+".tar.gz", func(files []os.DirEntry, i, j int) bool {
		ti, _ := getPackageTime(files[i].Name())
		tj, _ := getPackageTime(files[j].Name())
		return ti.Before(tj)
	})
	ex := extractor.NewTgz()
	dirName := newFile
	for filepath.Ext(dirName) != "" {
		dirName = strings.TrimSuffix(dirName, filepath.Ext(dirName))
	}
	dir := filepath.Join(m.Workspace, dirName)
	if _, err := os.Stat(dir); err == nil {
		return dirName, nil
	}
	logrus.Infof("extracting %s to %s...", filepath.Join(m.ArchiveDir, newFile), dir)
	if err := ex.Extract(filepath.Join(m.ArchiveDir, newFile), dir); err != nil {
		return dirName, err
	}
	if err := os.WriteFile(filepath.Join(dir, "__built__"), nil, 0640); err != nil {
		return dirName, err
	}
	return dirName, nil
}
