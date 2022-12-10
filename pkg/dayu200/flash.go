package dayu200

import (
	"fmt"
	"fotff/utils"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var partList = []string{"boot_linux", "system", "vendor", "userdata", "resource", "ramdisk", "chipset", "sys-prod", "chip-prod", "updater"}

func (m *Manager) flashDevice(pkg string) error {
	if err := m.tryRebootToLoader(); err != nil {
		return err
	}
	time.Sleep(5 * time.Second) // sleep a while for rebooting
	if err := m.flashImages(pkg); err != nil {
		return err
	}
	time.Sleep(5 * time.Second) // sleep a while for writing
	if err := utils.Exec(m.FlashTool, "RD"); err != nil {
		return fmt.Errorf("reboot device fail: %v", err)
	}
	time.Sleep(30 * time.Second) // sleep a while for rebooting
	logrus.Infof("flash device successfully")
	return nil
}

func (m *Manager) tryRebootToLoader() error {
	var hdc string
	if hdc, _ = exec.LookPath("hdc"); hdc == "" {
		hdc, _ = exec.LookPath("hdc_std")
	}
	if hdc == "" {
		logrus.Warn("can not find 'hdc', please install")
		return nil
	}
	for i := 0; i < 5; i++ {
		utils.Exec(hdc, "kill")
		time.Sleep(time.Second)
		utils.Exec(hdc, "start")
		time.Sleep(time.Second)
		out, err := utils.ExecCombinedOutput(hdc, "list", "targets")
		if err != nil {
			logrus.Errorf("failed to list hdc targets: %s", string(out))
			continue
		}
		lines := strings.Fields(string(out))
		for _, dev := range lines {
			if dev == "[Empty]" {
				logrus.Warn("can not find any hdc targets")
				break
			}
			if m.SN == "" {
				return utils.Exec(hdc, "shell", "reboot", "loader")
			}
			if dev == m.SN {
				return utils.Exec(hdc, "-t", m.SN, "shell", "reboot", "loader")
			}
		}
		logrus.Infof("%s not found", m.SN)
	}
	logrus.Warn("can not find any hdc targets, assume it has been in loader mode")
	return nil
}

func (m *Manager) flashImages(pkg string) error {
	logrus.Infof("calling flash tool for %s...", pkg)
	if err := utils.Exec(m.FlashTool, "UL", filepath.Join(pkg, "MiniLoaderAll.bin"), "-noreset"); err != nil {
		logrus.Errorf("flash MiniLoaderAll.bin fail: %v", err)
		return err
	}
	time.Sleep(3 * time.Second)
	if err := utils.Exec(m.FlashTool, "DI", "-p", filepath.Join(pkg, "parameter.txt")); err != nil {
		logrus.Errorf("flash parameter.txt fail: %v", err)
		return err
	}
	time.Sleep(5 * time.Second)
	if err := utils.Exec(m.FlashTool, "DI", "-uboot", filepath.Join(pkg, "uboot.img"), filepath.Join(pkg, "parameter.txt")); err != nil {
		logrus.Errorf("flash device fail: %v", err)
		return err
	}
	time.Sleep(5 * time.Second)
	for _, part := range partList {
		if _, err := os.Stat(filepath.Join(pkg, part+".img")); err != nil {
			if os.IsNotExist(err) {
				logrus.Infof("part %s.img not exist, ignored", part)
				continue
			}
			return err
		}
		if err := utils.Exec(m.FlashTool, "DI", "-"+part, filepath.Join(pkg, part+".img"), filepath.Join(pkg, "parameter.txt")); err != nil {
			logrus.Errorf("flash device fail: %v", err)
			logrus.Warnf("try again...")
			if err := utils.Exec(m.FlashTool, "DI", "-"+part, filepath.Join(pkg, part+".img"), filepath.Join(pkg, "parameter.txt")); err != nil {
				logrus.Errorf("flash device fail: %v", err)
				return err
			}
		}
		time.Sleep(3 * time.Second)
	}
	return nil
}
