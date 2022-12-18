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

package rec

import (
	"context"
	"fmt"
	"fotff/res"
	"fotff/tester"
	"github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

type FotffMocker struct {
	StepsNum   int
	FirstFail  int
	lock       sync.Mutex
	runningPkg map[string]string
}

func TestMain(m *testing.M) {
	defer os.RemoveAll(".fotff")
	defer os.RemoveAll("logs")
	m.Run()
}

func NewFotffMocker(stepsNum int, firstFail int) *FotffMocker {
	return &FotffMocker{
		StepsNum:   stepsNum,
		FirstFail:  firstFail,
		runningPkg: map[string]string{},
	}
}

func (f *FotffMocker) TaskName() string {
	return "mocker"
}

func (f *FotffMocker) DoTestTask(device string, ctx context.Context) ([]tester.Result, error) {
	time.Sleep(time.Duration(rand.Intn(1)) * time.Millisecond)
	select {
	case <-ctx.Done():
		return nil, context.Canceled
	default:
	}
	return []tester.Result{{TestCaseName: f.TestCaseName(), Status: tester.ResultFail}}, nil
}

func (f *FotffMocker) DoTestCase(device string, testcase string, ctx context.Context) (tester.Result, error) {
	time.Sleep(time.Duration(rand.Intn(1)) * time.Millisecond)
	select {
	case <-ctx.Done():
		return tester.Result{}, context.Canceled
	default:
	}
	f.lock.Lock()
	pkg, _ := strconv.Atoi(f.runningPkg[device])
	f.lock.Unlock()
	if pkg >= f.FirstFail {
		logrus.Infof("mock: test %s at %s done, result is %s", testcase, device, tester.ResultFail)
		return tester.Result{TestCaseName: testcase, Status: tester.ResultFail}, nil
	}
	logrus.Infof("mock: test %s at %s done, result is %s", testcase, device, tester.ResultPass)
	return tester.Result{TestCaseName: testcase, Status: tester.ResultPass}, nil
}

func (f *FotffMocker) Flash(device string, pkg string, ctx context.Context) error {
	time.Sleep(time.Duration(rand.Intn(1)) * time.Millisecond)
	select {
	case <-ctx.Done():
		return context.Canceled
	default:
	}
	f.lock.Lock()
	f.runningPkg[device] = pkg
	logrus.Infof("mock: flash %s to %s done", pkg, device)
	f.lock.Unlock()
	return nil
}

func (f *FotffMocker) LastIssue(pkg string) (string, error) {
	return "issue" + pkg, nil
}

func (f *FotffMocker) Steps(from, to string) (ret []string, err error) {
	fromIndex, _ := strconv.Atoi(from)
	toIndex, _ := strconv.Atoi(to)
	for i := fromIndex + 1; i <= toIndex; i++ {
		ret = append(ret, strconv.Itoa(i))
	}
	return
}

func (f *FotffMocker) GetNewer(cur string) (string, error) {
	return strconv.Itoa(f.StepsNum), nil
}

func (f *FotffMocker) TestCaseName() string {
	return "MOCK_FAILED_TEST_CASE"
}

func (f *FotffMocker) First() string {
	return "0"
}

func (f *FotffMocker) Last() string {
	return strconv.Itoa(f.StepsNum)
}

func TestFindOutTheFirstFail(t *testing.T) {
	tests := []struct {
		name   string
		mocker *FotffMocker
	}{
		{
			name:   "0-1(X)",
			mocker: NewFotffMocker(1, 1),
		},
		{
			name:   "0-1(X)-2",
			mocker: NewFotffMocker(2, 1),
		},
		{
			name:   "0-1-2(X)",
			mocker: NewFotffMocker(2, 2),
		},
		{
			name:   "0-1(X)-2-3",
			mocker: NewFotffMocker(3, 1),
		},
		{
			name:   "0-1-2(X)-3",
			mocker: NewFotffMocker(3, 2),
		},
		{
			name:   "0-1-2-3(X)",
			mocker: NewFotffMocker(3, 3),
		},
		{
			name:   "0-1(X)-2-3-4",
			mocker: NewFotffMocker(4, 1),
		},
		{
			name:   "0-1-2(X)-3-4",
			mocker: NewFotffMocker(4, 2),
		},
		{
			name:   "0-1-2-3(X)-4",
			mocker: NewFotffMocker(4, 3),
		},
		{
			name:   "0-1-2-3-4(X)",
			mocker: NewFotffMocker(4, 4),
		},
		{
			name:   "0-1(X)-2-3-4-5",
			mocker: NewFotffMocker(5, 1),
		},
		{
			name:   "0-1-2(X)-3-4-5",
			mocker: NewFotffMocker(5, 2),
		},
		{
			name:   "0-1-2-3(X)-4-5",
			mocker: NewFotffMocker(5, 3),
		},
		{
			name:   "0-1-2-3-4(X)-5",
			mocker: NewFotffMocker(5, 4),
		},
		{
			name:   "0-1-2-3-4-5(X)",
			mocker: NewFotffMocker(5, 5),
		},
		{
			name:   "0-1-2...262143(X)...1048575",
			mocker: NewFotffMocker(1048575, 262143),
		},
		{
			name:   "0-1-2...262144(X)...1048575",
			mocker: NewFotffMocker(1048575, 262144),
		},
		{
			name:   "0-1-2...262145(X)...1048575",
			mocker: NewFotffMocker(1048575, 262145),
		},
		{
			name:   "0-1-2...262143(X)...1048576",
			mocker: NewFotffMocker(1048576, 262143),
		},
		{
			name:   "0-1-2...262144(X)...1048576",
			mocker: NewFotffMocker(1048576, 262144),
		},
		{
			name:   "0-1-2...262145(X)...1048576",
			mocker: NewFotffMocker(1048576, 262145),
		},
		{
			name:   "0-1-2...262143(X)...1048577",
			mocker: NewFotffMocker(1048577, 262143),
		},
		{
			name:   "0-1-2...262144(X)...1048577",
			mocker: NewFotffMocker(1048577, 262144),
		},
		{
			name:   "0-1-2...262145(X)...1048577",
			mocker: NewFotffMocker(1048577, 262145),
		},
		{
			name:   "0-1-2...1234567(X)...10000000",
			mocker: NewFotffMocker(10000000, 1234567),
		},
		{
			name:   "0-1-2...1234567(X)...100000001",
			mocker: NewFotffMocker(10000001, 1234567),
		},
		{
			name:   "0-1-2...7654321(X)...10000000",
			mocker: NewFotffMocker(10000000, 7654321),
		},
		{
			name:   "0-1-2...7654321(X)...10000001",
			mocker: NewFotffMocker(10000001, 7654321),
		},
		{
			name:   "0-1(X)-2...10000000",
			mocker: NewFotffMocker(10000000, 1),
		},
		{
			name:   "0-1(X)-2...10000001",
			mocker: NewFotffMocker(10000001, 1),
		},
		{
			name:   "0-1-2...10000000(X)",
			mocker: NewFotffMocker(10000000, 10000000),
		},
		{
			name:   "0-1-2...10000001(X)",
			mocker: NewFotffMocker(10000001, 10000001),
		},
	}
	for i := 1; i <= 5; i++ {
		res.Fake(i)
		for _, tt := range tests {
			t.Run(fmt.Sprintf("RES%d:%s", i, tt.name), func(t *testing.T) {
				ret, err := FindOutTheFirstFail(tt.mocker, tt.mocker, tt.mocker.TestCaseName(), tt.mocker.First(), tt.mocker.Last())
				if err != nil {
					t.Errorf("err: expcect: <nil>, actual: %v", err)
				}
				expectIssue, _ := tt.mocker.LastIssue(strconv.Itoa(tt.mocker.FirstFail))
				if ret != expectIssue {
					t.Errorf("fotff result: expect: %s, actual: %s", expectIssue, ret)
				}
			})
		}
	}
}
