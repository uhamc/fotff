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

type Tester struct{}

func NewTester() tester.Tester {
	return &Tester{}
}

func (t *Tester) TaskName() string {
	return "mock"
}

func (t *Tester) DoTestTask(device string, ctx context.Context) ([]tester.Result, error) {
	logrus.Infof("TEST_001 pass")
	logrus.Warnf("TEST_002 pass")
	logrus.Warnf("TEST_003 pass")
	return []tester.Result{
		{TestCaseName: "TEST_001", Status: tester.ResultPass},
		{TestCaseName: "TEST_002", Status: tester.ResultPass},
		{TestCaseName: "TEST_003", Status: tester.ResultPass},
	}, nil
}

func (t *Tester) DoTestCase(device string, testCase string, ctx context.Context) (tester.Result, error) {
	logrus.Warnf("%s pass", testCase)
	return tester.Result{TestCaseName: testCase, Status: tester.ResultPass}, nil
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
