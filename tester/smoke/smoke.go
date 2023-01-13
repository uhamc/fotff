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

package smoke

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"fotff/tester"
	"fotff/utils"
	"github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Tester struct {
	Py         string `key:"py"`
	Config     string `key:"config"`
	AnswerPath string `key:"answer_path"`
	SavePath   string `key:"save_path"`
	ToolsPath  string `key:"tools_path"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewTester() tester.Tester {
	ret := &Tester{}
	utils.ParseFromConfigFile("smoke", ret)
	return ret
}

func (t *Tester) TaskName() string {
	return "smoke_test"
}

func (t *Tester) DoTestTask(deviceSN string, ctx context.Context) (ret []tester.Result, err error) {
	reportDir := fmt.Sprintf("%X", md5.Sum([]byte(fmt.Sprintf("%d", rand.Int()))))
	if err := os.MkdirAll(filepath.Join(t.SavePath, reportDir), 0755); err != nil {
		return nil, err
	}
	args := []string{t.Py, "--config", t.Config, "--answer_path", t.AnswerPath, "--save_path", filepath.Join(t.SavePath, reportDir), "--tools_path", t.ToolsPath}
	if deviceSN != "" {
		args = append(args, "--device_num", deviceSN)
	}
	if err := utils.ExecContext(ctx, "python", args...); err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, err
		}
		logrus.Errorf("do test suite fail: %v", err)
		return nil, err
	}
	return t.readReport(reportDir)
}

func (t *Tester) DoTestCase(deviceSN, testCase string, ctx context.Context) (ret tester.Result, err error) {
	reportDir := fmt.Sprintf("%X", md5.Sum([]byte(fmt.Sprintf("%d", rand.Int()))))
	if err := os.MkdirAll(filepath.Join(t.SavePath, reportDir), 0755); err != nil {
		return ret, err
	}
	args := []string{t.Py, "--config", t.Config, "--answer_path", t.AnswerPath, "--save_path", filepath.Join(t.SavePath, reportDir), "--tools_path", t.ToolsPath, "--test_num", testCase}
	if deviceSN != "" {
		args = append(args, "--device_num", deviceSN)
	}
	if err := utils.ExecContext(ctx, "python", args...); err != nil {
		if errors.Is(err, context.Canceled) {
			return ret, err
		}
		logrus.Errorf("do test case %s fail: %v", testCase, err)
		return ret, err
	}
	r, err := t.readReport(reportDir)
	if len(r) == 0 {
		return ret, fmt.Errorf("read latest report err, no result found")
	}
	if r[0].TestCaseName != testCase {
		return ret, fmt.Errorf("read latest report err, no matched result found")
	}
	logrus.Infof("do testcase %s at %s done, result is %s", r[0].TestCaseName, deviceSN, r[0].Status)
	return r[0], nil
}

func (t *Tester) DoTestCases(deviceSN string, testcases []string, ctx context.Context) (ret []tester.Result, err error) {
	reportDir := fmt.Sprintf("%X", md5.Sum([]byte(fmt.Sprintf("%d", rand.Int()))))
	if err := os.MkdirAll(filepath.Join(t.SavePath, reportDir), 0755); err != nil {
		return nil, err
	}
	args := []string{t.Py, "--config", t.Config, "--answer_path", t.AnswerPath, "--save_path", filepath.Join(t.SavePath, reportDir), "--tools_path", t.ToolsPath, "--test_num", strings.Join(testcases, " ")}
	if deviceSN != "" {
		args = append(args, "--device_num", deviceSN)
	}
	if err := utils.ExecContext(ctx, "python", args...); err != nil {
		if errors.Is(err, context.Canceled) {
			return ret, err
		}
		logrus.Errorf("do test cases %v fail: %v", testcases, err)
		return ret, err
	}
	return t.readReport(reportDir)
}

func (t *Tester) readReport(reportDir string) (ret []tester.Result, err error) {
	data, err := os.ReadFile(filepath.Join(t.SavePath, reportDir, "result.json"))
	if err != nil {
		logrus.Errorf("read report json fail: %v", err)
		return nil, err
	}
	var result []struct {
		TestCaseName int    `json:"test_case_name"`
		Status       string `json:"status"`
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		logrus.Errorf("unmarshal report xml fail: %v", err)
		return nil, err
	}
	for _, r := range result {
		if r.Status == "pass" {
			ret = append(ret, tester.Result{TestCaseName: strconv.Itoa(r.TestCaseName), Status: tester.ResultPass})
		} else {
			ret = append(ret, tester.Result{TestCaseName: strconv.Itoa(r.TestCaseName), Status: tester.ResultFail})
		}
	}
	return ret, err
}
