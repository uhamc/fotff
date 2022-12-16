package dayu200

import (
	"context"
	"errors"
	"fmt"
	"fotff/utils"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const enableTestModeScript = `mount -o rw,remount /; param set persist.ace.testmode.enabled 1; param set persist.sys.hilog.debug.on true; sed -i 's/enforcing/permissive/g' /system/etc/selinux/config; sync; reboot`

var partList = []string{"boot_linux", "system", "vendor", "userdata", "resource", "ramdisk", "chipset", "sys-prod", "chip-prod", "updater"}

func (m *Manager) flashDevice(device string, pkg string, ctx context.Context) error {
	if err := m.tryRebootToLoader(device); err != nil {
		return err
	}
	if err := m.flashImages(device, pkg); err != nil {
		return err
	}
	time.Sleep(20 * time.Second) // usually, it takes about 20s to reboot into OpenHarmony
	if err := m.enableTestMode(device); err != nil {
		return err
	}
	time.Sleep(10 * time.Second) // wait 10s more to ensure system has been started completely
	logrus.Infof("flash device successfully")
	return nil
}

func (m *Manager) tryRebootToLoader(device string) error {
	logrus.Info("try to reboot to loader...")
	defer time.Sleep(5 * time.Second) // sleep a while for rebooting to loader
	if connected := m.waitHDC(device, 20*time.Second); connected {
		if device == "" {
			return utils.Exec(m.hdc, "shell", "reboot", "loader")
		} else {
			return utils.Exec(m.hdc, "-t", device, "shell", "reboot", "loader")
		}
	}
	logrus.Warn("can not find any hdc targets, assume it has been in loader mode")
	return nil
}

func (m *Manager) waitHDC(device string, timeout time.Duration) bool {
	timer := time.NewTimer(timeout)
	for {
		select {
		case <-timer.C:
			return false
		default:
		}
		utils.Exec(m.hdc, "kill")
		time.Sleep(time.Second)
		utils.Exec(m.hdc, "start")
		time.Sleep(time.Second)
		out, err := utils.ExecCombinedOutput(m.hdc, "list", "targets")
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
			if device == "" || dev == device {
				return true
			}
		}
		logrus.Infof("%s not found", device)
	}
}

func (m *Manager) flashImages(device string, pkg string) error {
	logrus.Infof("calling flash tool for %s...", pkg)
	locationID := m.locations[device]
	if locationID == "" {
		data, _ := utils.ExecCombinedOutput(m.FlashTool, "LD")
		locationID = strings.TrimPrefix(regexp.MustCompile(`LocationID=\d+`).FindString(string(data)), "LocationID=")
		if locationID == "" {
			time.Sleep(5 * time.Second)
			data, _ := utils.ExecCombinedOutput(m.FlashTool, "LD")
			locationID = strings.TrimPrefix(regexp.MustCompile(`LocationID=\d+`).FindString(string(data)), "LocationID=")
		}
	}
	if err := utils.Exec(m.FlashTool, "-s", locationID, "UL", filepath.Join(m.Workspace, pkg, "MiniLoaderAll.bin"), "-noreset"); err != nil {
		logrus.Errorf("flash MiniLoaderAll.bin fail: %v", err)
		time.Sleep(5 * time.Second)
		if err := utils.Exec(m.FlashTool, "-s", locationID, "UL", filepath.Join(m.Workspace, pkg, "MiniLoaderAll.bin"), "-noreset"); err != nil {
			logrus.Errorf("flash MiniLoaderAll.bin fail: %v", err)
			return err
		}
	}
	time.Sleep(3 * time.Second)
	if err := utils.Exec(m.FlashTool, "-s", locationID, "DI", "-p", filepath.Join(m.Workspace, pkg, "parameter.txt")); err != nil {
		logrus.Errorf("flash parameter.txt fail: %v", err)
		return err
	}
	time.Sleep(5 * time.Second)
	if err := utils.Exec(m.FlashTool, "-s", locationID, "DI", "-uboot", filepath.Join(m.Workspace, pkg, "uboot.img"), filepath.Join(m.Workspace, pkg, "parameter.txt")); err != nil {
		logrus.Errorf("flash device fail: %v", err)
		return err
	}
	time.Sleep(5 * time.Second)
	for _, part := range partList {
		if _, err := os.Stat(filepath.Join(m.Workspace, pkg, part+".img")); err != nil {
			if os.IsNotExist(err) {
				logrus.Infof("part %s.img not exist, ignored", part)
				continue
			}
			return err
		}
		if err := utils.Exec(m.FlashTool, "-s", locationID, "DI", "-"+part, filepath.Join(m.Workspace, pkg, part+".img"), filepath.Join(m.Workspace, pkg, "parameter.txt")); err != nil {
			logrus.Errorf("flash device fail: %v", err)
			logrus.Warnf("try again...")
			if err := utils.Exec(m.FlashTool, "-s", locationID, "DI", "-"+part, filepath.Join(m.Workspace, pkg, part+".img"), filepath.Join(m.Workspace, pkg, "parameter.txt")); err != nil {
				logrus.Errorf("flash device fail: %v", err)
				return err
			}
		}
		time.Sleep(3 * time.Second)
	}
	time.Sleep(5 * time.Second) // sleep a while for writing
	if err := utils.Exec(m.FlashTool, "-s", locationID, "RD"); err != nil {
		return fmt.Errorf("reboot device fail: %v", err)
	}
	return nil
}

func (m *Manager) enableTestMode(device string) (err error) {
	if connected := m.waitHDC(device, time.Minute); !connected {
		return errors.New("can not connect to hdc, timeout")
	}
	logrus.Info("try to enable test mode...")
	if device == "" {
		err = utils.Exec(m.hdc, "shell", enableTestModeScript)
	} else {
		err = utils.Exec(m.hdc, "-t", device, "shell", enableTestModeScript)
	}
	if err != nil {
		return err
	}
	time.Sleep(20 * time.Second) // usually, it takes about 20s to reboot into OpenHarmony
	if connected := m.waitHDC(device, time.Minute); !connected {
		return errors.New("can not connect to hdc, timeout")
	}
	return nil
}
