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
	"fotff/tester"
	"github.com/sirupsen/logrus"
)

type Tester struct {
	pkgCount int
}

func NewTester() tester.Tester {
	return &Tester{pkgCount: -1}
}

func (t *Tester) TaskName() string {
	return "mock"
}

func (t *Tester) DoTestTask(device string, ctx context.Context) ([]tester.Result, error) {
	t.pkgCount++
	if t.pkgCount%2 == 0 {
		logrus.Infof("TEST_001 pass")
		logrus.Infof("TEST_002 pass")
		return []tester.Result{
			{TestCaseName: "TEST_001", Status: tester.ResultPass},
			{TestCaseName: "TEST_002", Status: tester.ResultPass},
		}, nil
	}
	logrus.Infof("TEST_001 pass")
	logrus.Warnf("TEST_002 fail")
	return []tester.Result{
		{TestCaseName: "TEST_001", Status: tester.ResultPass},
		{TestCaseName: "TEST_002", Status: tester.ResultFail},
	}, nil
}

func (t *Tester) DoTestCase(device string, testCase string, ctx context.Context) (tester.Result, error) {
	if t.pkgCount%2 == 0 {
		logrus.Infof("%s pass", testCase)
		return tester.Result{TestCaseName: testCase, Status: tester.ResultPass}, nil
	}
	switch testCase {
	case "TEST_001":
		logrus.Infof("%s pass", testCase)
		return tester.Result{TestCaseName: testCase, Status: tester.ResultPass}, nil
	case "TEST_002":
		logrus.Warnf("%s fail", testCase)
		return tester.Result{TestCaseName: testCase, Status: tester.ResultFail}, nil
	default:
		panic("not defined")
	}
}

func (t *Tester) DoTestCases(device string, testcases []string, ctx context.Context) ([]tester.Result, error) {
	var ret []tester.Result
	for _, testcase := range testcases {
		r, err := t.DoTestCase(device, testcase, ctx)
		if err != nil {
			return nil, err
		}
		ret = append(ret, r)
	}
	return ret, nil
}
