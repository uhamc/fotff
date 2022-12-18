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

package mock

import (
	"context"
	"fmt"
	"fotff/pkg"
	"github.com/sirupsen/logrus"
	"time"
)

type Manager struct {
	pkgCount int
}

func NewManager() pkg.Manager {
	return &Manager{}
}

func (m *Manager) LastIssue(pkg string) (string, error) {
	ret := fmt.Sprintf("https://testserver.com/issues/%s", pkg)
	logrus.Infof("LastIssue: mock implementation returns %s", ret)
	return ret, nil
}

func (m *Manager) Steps(from, to string) ([]string, error) {
	var ret = []string{"step1", "step2", "step3"}
	for i := range ret {
		ret[i] = fmt.Sprintf("%s-%s-%s", from, to, ret[i])
	}
	logrus.Infof("Steps: mock implementation returns %v", ret)
	return ret, nil
}

func (m *Manager) GetNewer(cur string) (string, error) {
	ret := fmt.Sprintf("pkg%d", m.pkgCount)
	time.Sleep(time.Duration(m.pkgCount) * time.Second)
	m.pkgCount++
	logrus.Infof("GetNewer: mock implementation returns %s", ret)
	return ret, nil
}

func (m *Manager) Flash(device string, pkg string, ctx context.Context) error {
	time.Sleep(time.Second)
	logrus.Infof("Flash: flashing %s to %s, mock implementation returns OK unconditionally", pkg, device)
	return nil
}
