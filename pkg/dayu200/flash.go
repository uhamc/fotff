/*
 * Copyright (c) 2022 Huawei Device Co., Ltd.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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

// All timeouts are calculated on normal cases, we do not certain that timeouts are enough if some sleeps canceled.
// So simply we do not cancel any Sleep(). TODO: use utils.SleepContext() instead.
func (m *Manager) flashDevice(device string, pkg string, ctx context.Context) error {
	if err := m.tryRebootToLoader(device, ctx); err != nil {
		return err
	}
	if err := m.flashImages(device, pkg, ctx); err != nil {
		return err
	}
	time.Sleep(20 * time.Second) // usually, it takes about 20s to reboot into OpenHarmony
	if err := m.enableTestMode(device, ctx); err != nil {
		return err
	}
	time.Sleep(10 * time.Second) // wait 10s more to ensure system has been started completely
	logrus.Infof("flash device %s successfully", device)
	return nil
}

func (m *Manager) tryRebootToLoader(device string, ctx context.Context) error {
	logrus.Infof("try to reboot %s to loader...", device)
	defer time.Sleep(5 * time.Second) // sleep a while for rebooting to loader
	waitCtx, cancelFn := context.WithTimeout(ctx, 20*time.Second)
	defer cancelFn()
	if connected := m.waitHDC(device, waitCtx); connected {
		if device == "" {
			return utils.ExecContext(ctx, m.hdc, "shell", "reboot", "loader")
		} else {
			return utils.ExecContext(ctx, m.hdc, "-t", device, "shell", "reboot", "loader")
		}
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	logrus.Warn("can not find target hdc device, assume it has been in loader mode")
	return nil
}

func (m *Manager) waitHDC(device string, ctx context.Context) bool {
	for {
		select {
		case <-ctx.Done():
			return false
		default:
		}
		utils.ExecContext(ctx, m.hdc, "kill")
		time.Sleep(time.Second)
		utils.ExecContext(ctx, m.hdc, "start")
		time.Sleep(time.Second)
		out, err := utils.ExecCombinedOutputContext(ctx, m.hdc, "list", "targets")
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return false
			}
			logrus.Errorf("failed to list hdc targets: %s, %s", string(out), err)
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

func (m *Manager) flashImages(device string, pkg string, ctx context.Context) error {
	logrus.Infof("calling flash tool to flash %s into %s...", pkg, device)
	locationID := m.locations[device]
	if locationID == "" {
		data, _ := utils.ExecCombinedOutputContext(ctx, m.FlashTool, "LD")
		locationID = strings.TrimPrefix(regexp.MustCompile(`LocationID=\d+`).FindString(string(data)), "LocationID=")
		if locationID == "" {
			time.Sleep(5 * time.Second)
			data, _ := utils.ExecCombinedOutputContext(ctx, m.FlashTool, "LD")
			locationID = strings.TrimPrefix(regexp.MustCompile(`LocationID=\d+`).FindString(string(data)), "LocationID=")
		}
	}
	logrus.Infof("locationID of %s is [%s]", device, locationID)
	if err := utils.ExecContext(ctx, m.FlashTool, "-s", locationID, "UL", filepath.Join(m.Workspace, pkg, "MiniLoaderAll.bin"), "-noreset"); err != nil {
		if errors.Is(err, context.Canceled) {
			return err
		}
		logrus.Errorf("flash MiniLoaderAll.bin fail: %v", err)
		time.Sleep(5 * time.Second)
		if err := utils.ExecContext(ctx, m.FlashTool, "-s", locationID, "UL", filepath.Join(m.Workspace, pkg, "MiniLoaderAll.bin"), "-noreset"); err != nil {
			if errors.Is(err, context.Canceled) {
				return err
			}
			logrus.Errorf("flash MiniLoaderAll.bin fail: %v", err)
			return err
		}
	}
	time.Sleep(3 * time.Second)
	if err := utils.ExecContext(ctx, m.FlashTool, "-s", locationID, "DI", "-p", filepath.Join(m.Workspace, pkg, "parameter.txt")); err != nil {
		if errors.Is(err, context.Canceled) {
			return err
		}
		logrus.Errorf("flash parameter.txt fail: %v", err)
		return err
	}
	time.Sleep(5 * time.Second)
	if err := utils.ExecContext(ctx, m.FlashTool, "-s", locationID, "DI", "-uboot", filepath.Join(m.Workspace, pkg, "uboot.img"), filepath.Join(m.Workspace, pkg, "parameter.txt")); err != nil {
		if errors.Is(err, context.Canceled) {
			return err
		}
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
		if err := utils.ExecContext(ctx, m.FlashTool, "-s", locationID, "DI", "-"+part, filepath.Join(m.Workspace, pkg, part+".img"), filepath.Join(m.Workspace, pkg, "parameter.txt")); err != nil {
			if errors.Is(err, context.Canceled) {
				return err
			}
			logrus.Errorf("flash device fail: %v", err)
			logrus.Warnf("try again...")
			if err := utils.ExecContext(ctx, m.FlashTool, "-s", locationID, "DI", "-"+part, filepath.Join(m.Workspace, pkg, part+".img"), filepath.Join(m.Workspace, pkg, "parameter.txt")); err != nil {
				if errors.Is(err, context.Canceled) {
					return err
				}
				logrus.Errorf("flash device fail: %v", err)
				return err
			}
		}
		time.Sleep(3 * time.Second)
	}
	time.Sleep(5 * time.Second) // sleep a while for writing
	if err := utils.ExecContext(ctx, m.FlashTool, "-s", locationID, "RD"); err != nil {
		if errors.Is(err, context.Canceled) {
			return err
		}
		return fmt.Errorf("reboot device fail: %v", err)
	}
	return nil
}

func (m *Manager) enableTestMode(device string, ctx context.Context) (err error) {
	waitCtx1, cancelFn1 := context.WithTimeout(ctx, time.Minute)
	defer cancelFn1()
	if connected := m.waitHDC(device, waitCtx1); !connected {
		if err := ctx.Err(); err != nil {
			return err
		}
		return fmt.Errorf("can not connect %s to hdc, timeout", device)
	}
	logrus.Info("try to enable test mode...")
	if device == "" {
		err = utils.ExecContext(ctx, m.hdc, "shell", enableTestModeScript)
	} else {
		err = utils.ExecContext(ctx, m.hdc, "-t", device, "shell", enableTestModeScript)
	}
	if err != nil {
		return err
	}
	time.Sleep(20 * time.Second) // usually, it takes about 20s to reboot into OpenHarmony
	waitCtx2, cancelFn2 := context.WithTimeout(ctx, time.Minute)
	defer cancelFn2()
	if connected := m.waitHDC(device, waitCtx2); !connected {
		if err := ctx.Err(); err != nil {
			return err
		}
		return fmt.Errorf("can not connect %s to hdc, timeout", device)
	}
	return nil
}
