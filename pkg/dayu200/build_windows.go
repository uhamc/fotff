//go:build windows

package dayu200

import (
	"fmt"
	"fotff/utils"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

func (m *Manager) build(pkg string) error {
	if _, err := os.Stat(filepath.Join(pkg, "__built__")); err == nil {
		return nil
	}
	cmd := fmt.Sprintf("mkdir -p %s && cd %s && repo init -u https://gitee.com/openharmony/manifest.git", m.BuildServerConfig.BuildWorkSpace, m.BuildServerConfig.BuildWorkSpace)
	if out, err := utils.RunCmdViaSSH(m.BuildServerConfig.Addr, m.BuildServerConfig.User, m.BuildServerConfig.Passwd, cmd); err != nil {
		logrus.Error(out)
		return fmt.Errorf("remote: mkdir error: %s", err)
	}
	if err := utils.TransFileViaSSH(utils.Upload, m.BuildServerConfig.Addr, m.BuildServerConfig.User, m.BuildServerConfig.Passwd,
		fmt.Sprintf("%s/.repo/manifest.xml", m.BuildServerConfig.BuildWorkSpace), filepath.Join(pkg, "manifest_tag.xml")); err != nil {
		return fmt.Errorf("upload and replace manifest error: %s", err)
	}
	cmd = fmt.Sprintf("cd %s && repo sync -c --no-tags --force-remove-dirty", m.BuildServerConfig.BuildWorkSpace)
	if out, err := utils.RunCmdViaSSH(m.BuildServerConfig.Addr, m.BuildServerConfig.User, m.BuildServerConfig.Passwd, cmd); err != nil {
		logrus.Error(out)
		return fmt.Errorf("remote: repo sync error: %s", err)
	}
	cmd = fmt.Sprintf("cd %s && %s", m.BuildServerConfig.BuildWorkSpace, preCompileCMD)
	if out, err := utils.RunCmdViaSSH(m.BuildServerConfig.Addr, m.BuildServerConfig.User, m.BuildServerConfig.Passwd, cmd); err != nil {
		logrus.Error(out)
		return fmt.Errorf("remote: pre-compile command error: %s", err)
	}
	cmd = fmt.Sprintf("cd %s && %s", m.BuildServerConfig.BuildWorkSpace, compileCMD)
	if out, err := utils.RunCmdViaSSH(m.BuildServerConfig.Addr, m.BuildServerConfig.User, m.BuildServerConfig.Passwd, cmd); err != nil {
		logrus.Error(out)
		return fmt.Errorf("remote: compile command error: %s", err)
	}
	for _, f := range imgList {
		imgName := filepath.Base(f)
		if err := utils.TransFileViaSSH(utils.Download, m.BuildServerConfig.Addr, m.BuildServerConfig.User, m.BuildServerConfig.Passwd,
			fmt.Sprintf("%s/%s", m.BuildServerConfig.BuildWorkSpace, f), filepath.Join(pkg, imgName)); err != nil {
			return fmt.Errorf("download file %s error: %s", f, err)
		}
	}
	return nil
}
