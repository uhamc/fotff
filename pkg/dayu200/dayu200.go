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
	"code.cloudfoundry.org/archiver/extractor"
	"context"
	"fmt"
	"fotff/pkg"
	"fotff/res"
	"fotff/utils"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Manager struct {
	ArchiveDir     string `key:"archive_dir" default:"."`
	FromCI         string `key:"download_from_ci" default:"false"`
	Workspace      string `key:"workspace" default:"."`
	Branch         string `key:"branch" default:"master"`
	FlashTool      string `key:"flash_tool" default:"python"`
	LocationIDList string `key:"location_id_list"`
	locations      map[string]string
	fromCI         bool
	hdc            string
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
	var err error
	if ret.fromCI, err = strconv.ParseBool(ret.FromCI); err != nil {
		logrus.Panicf("can not parse 'download_from_ci', please check")
	}
	devs := res.DeviceList()
	locs := strings.Split(ret.LocationIDList, ",")
	if len(devs) != len(locs) {
		logrus.Panicf("location_id_list and devices mismatch")
	}
	ret.locations = map[string]string{}
	for i, loc := range locs {
		ret.locations[devs[i]] = loc
	}
	go ret.cleanupOutdated()
	return &ret
}

func (m *Manager) cleanupOutdated() {
	t := time.NewTicker(24 * time.Hour)
	for {
		<-t.C
		es, err := os.ReadDir(m.Workspace)
		if err != nil {
			logrus.Errorf("can not read %s: %v", m.Workspace, err)
			continue
		}
		for _, e := range es {
			if !e.IsDir() {
				continue
			}
			path := filepath.Join(m.Workspace, e.Name())
			info, err := e.Info()
			if err != nil {
				logrus.Errorf("can not read %s info: %v", path, err)
				continue
			}
			if time.Now().Sub(info.ModTime()) > 7*24*time.Hour {
				logrus.Warnf("%s outdated, cleanning up its contents...", path)
				m.cleanupPkgFiles(path)
			}
		}
	}
}

func (m *Manager) cleanupPkgFiles(path string) {
	es, err := os.ReadDir(path)
	if err != nil {
		logrus.Errorf("can not read %s: %v", path, err)
		return
	}
	for _, e := range es {
		if e.Name() == "manifest_tag.xml" || e.Name() == "__last_issue__" {
			continue
		}
		if err := os.RemoveAll(filepath.Join(path, e.Name())); err != nil {
			logrus.Errorf("remove %s fail: %v", filepath.Join(path, e.Name()), err)
		}
	}
}

func (m *Manager) Flash(device string, pkg string, ctx context.Context) error {
	logrus.Infof("now flash %s", pkg)
	if _, err := os.Stat(filepath.Join(m.Workspace, pkg, "__built__")); err != nil {
		if err := m.build(pkg, ctx); err != nil {
			logrus.Errorf("build pkg %s err: %v", pkg, err)
			return err
		}
	}
	return m.flashDevice(device, pkg, ctx)
}

func (m *Manager) Steps(from, to string) (pkgs []string, err error) {
	if from == to {
		return nil, fmt.Errorf("steps err: 'from' %s and 'to' %s are the same", from, to)
	}
	if c, found := utils.CacheGet("dayu200_steps", from+"__to__"+to); found {
		logrus.Infof("steps from %s to %s are cached", from, to)
		logrus.Infof("steps: %v", c.([]string))
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
	var newFile string
	if m.fromCI {
		newFile = m.getNewerFromCI(cur + ".tar.gz")
	} else {
		newFile = pkg.GetNewerFileFromDir(m.ArchiveDir, cur+".tar.gz", func(files []os.DirEntry, i, j int) bool {
			ti, _ := getPackageTime(files[i].Name())
			tj, _ := getPackageTime(files[j].Name())
			return ti.Before(tj)
		})
	}
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
