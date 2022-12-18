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
	"fmt"
	"fotff/res"
	"fotff/utils"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

func (m *Manager) build(pkg string, ctx context.Context) error {
	logrus.Infof("now build %s", pkg)
	server := res.GetBuildServer()
	defer res.ReleaseBuildServer(server)
	if _, err := os.Stat(filepath.Join(m.Workspace, pkg, "__built__")); err == nil {
		return nil
	}
	cmd := fmt.Sprintf("mkdir -p %s && cd %s && repo init -u https://gitee.com/openharmony/manifest.git", server.WorkSpace, server.WorkSpace)
	if err := utils.RunCmdViaSSHContext(ctx, server.Addr, server.User, server.Passwd, cmd); err != nil {
		return fmt.Errorf("remote: mkdir error: %w", err)
	}
	if err := utils.TransFileViaSSH(utils.Upload, server.Addr, server.User, server.Passwd,
		fmt.Sprintf("%s/.repo/manifest.xml", server.WorkSpace), filepath.Join(m.Workspace, pkg, "manifest_tag.xml")); err != nil {
		return fmt.Errorf("upload and replace manifest error: %w", err)
	}
	cmd = fmt.Sprintf("cd %s && repo sync -c --no-tags --force-remove-dirty && repo forall -c 'git reset --hard && git clean -dfx && git lfs update --force && git lfs install && git lfs pull'", server.WorkSpace)
	if err := utils.RunCmdViaSSHContext(ctx, server.Addr, server.User, server.Passwd, cmd); err != nil {
		return fmt.Errorf("remote: repo sync error: %w", err)
	}
	cmd = fmt.Sprintf("cd %s && %s", server.WorkSpace, preCompileCMD)
	if err := utils.RunCmdViaSSHContext(ctx, server.Addr, server.User, server.Passwd, cmd); err != nil {
		return fmt.Errorf("remote: pre-compile command error: %w", err)
	}
	cmd = fmt.Sprintf("cd %s && %s", server.WorkSpace, compileCMD)
	if err := utils.RunCmdViaSSHContext(ctx, server.Addr, server.User, server.Passwd, cmd); err != nil {
		return fmt.Errorf("remote: compile command error: %w", err)
	}
	// build already, pitiful if canceled, so continue copying
	for _, f := range imgList {
		imgName := filepath.Base(f)
		if err := utils.TransFileViaSSH(utils.Download, server.Addr, server.User, server.Passwd,
			fmt.Sprintf("%s/%s", server.WorkSpace, f), filepath.Join(m.Workspace, pkg, imgName)); err != nil {
			return fmt.Errorf("download file %s error: %w", f, err)
		}
	}
	return os.WriteFile(filepath.Join(m.Workspace, pkg, "__built__"), nil, 0640)
}
