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
	"errors"
	"fmt"
	"fotff/pkg"
	"fotff/res"
	"fotff/tester"
	"fotff/utils"
	"github.com/sirupsen/logrus"
	"math"
	"sync"
)

type cancelCtx struct {
	ctx context.Context
	fn  context.CancelFunc
}

// FindOutTheFirstFail returns the first issue URL that introduce the failure.
// 'fellows' are optional, these testcases may be tested with target testcase together.
func FindOutTheFirstFail(m pkg.Manager, t tester.Tester, testCase string, successPkg string, failPkg string, fellows ...string) (string, error) {
	if successPkg == "" {
		return "", fmt.Errorf("can not get a success package for %s", testCase)
	}
	steps, err := m.Steps(successPkg, failPkg)
	if err != nil {
		return "", err
	}
	return findOutTheFirstFail(m, t, testCase, steps, fellows...)
}

// findOutTheFirstFail is the recursive implementation to find out the first issue URL that introduce the failure.
// Arg steps' length must be grater than 1. The last step is a pre-known failure, while the rests are not tested.
// 'fellows' are optional. In the last recursive term, they have the same result as what the target testcases has.
// These fellows can be tested with target testcase together in this term to accelerate testing.
func findOutTheFirstFail(m pkg.Manager, t tester.Tester, testcase string, steps []string, fellows ...string) (string, error) {
	if len(steps) == 0 {
		return "", errors.New("steps are no between (success, failure]")
	}
	logrus.Infof("now use %d-section search to find out the first fault, the length of range is %d, between [%s, %s]", res.Num()+1, len(steps), steps[0], steps[len(steps)-1])
	if len(steps) == 1 {
		return m.LastIssue(steps[0])
	}
	// calculate gaps between every check point of N-section search. At least 1, or will cause duplicated tests.
	gapLen := float64(len(steps)-1) / float64(res.Num()+1)
	if gapLen < 1 {
		gapLen = 1
	}
	// 'success' and 'fail' record the left/right steps indexes of the next term recursive call.
	// Here defines functions and surrounding helpers to update success/fail indexes and cancel un-needed tests.
	success, fail := -1, len(steps)-1
	var lock sync.Mutex
	var contexts []cancelCtx
	updateRange := func(pass bool, index int) {
		lock.Lock()
		defer lock.Unlock()
		if pass && index > success {
			success = index
			for _, ctx := range contexts {
				if ctx.ctx.Value("index").(int) < success {
					ctx.fn()
				}
			}
		}
		if !pass && index < fail {
			fail = index
			for _, ctx := range contexts {
				if ctx.ctx.Value("index").(int) > fail {
					ctx.fn()
				}
			}
		}
	}
	// Now, start all tests concurrently.
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 1; i <= res.Num(); i++ {
		// Since the last step is a pre-known failure, we start index from the tail to avoid testing the last one.
		// Otherwise, if the last step is the only one we test this term, we can not narrow ranges to continue.
		index := len(steps) - 1 - int(math.Round(float64(i)*gapLen))
		if index < 0 {
			break
		}
		ctx, fn := context.WithCancel(context.WithValue(context.TODO(), "index", index))
		contexts = append(contexts, cancelCtx{ctx: ctx, fn: fn})
		wg.Add(1)
		go func(index int, ctx context.Context) {
			defer wg.Done()
			// Start after all test goroutine's contexts are registered.
			// Otherwise, contexts that not registered yet may out of controlling.
			<-start
			var pass bool
			var err error
			pass, fellows, err = flashAndTest(m, t, steps[index], testcase, ctx, fellows...)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					logrus.Warnf("abort to flash %s and test %s: %v", steps[index], testcase, err)
				} else {
					logrus.Errorf("flash %s and test %s fail: %v", steps[index], testcase, err)
				}
				return
			}
			updateRange(pass, index)
		}(index, ctx)
	}
	close(start)
	wg.Wait()
	if fail-success == len(steps) {
		return "", errors.New("all judgements failed, can not narrow ranges to continue")
	}
	return findOutTheFirstFail(m, t, testcase, steps[success+1:fail+1], fellows...)
}

func flashAndTest(m pkg.Manager, t tester.Tester, pkg string, testcase string, ctx context.Context, fellows ...string) (bool, []string, error) {
	var newFellows []string
	if status, found := utils.CacheGet("testcase_result", testcase+"__at__"+pkg); found {
		for _, fellow := range fellows {
			if fellowStatus, fellowFound := utils.CacheGet("testcase_result", fellow+"__at__"+pkg); fellowFound {
				if fellowStatus.(tester.Result).Status == status.(tester.Result).Status {
					newFellows = append(newFellows, fellow)
				}
			}
		}
		return status.(tester.Result).Status == tester.ResultPass, newFellows, nil
	}
	device := res.GetDevice()
	defer res.ReleaseDevice(device)
	if err := m.Flash(device, pkg, ctx); err != nil {
		return false, newFellows, err
	}
	results, err := t.DoTestCases(device, append(fellows, testcase), ctx)
	if err != nil {
		return false, newFellows, err
	}
	var testcaseStatus tester.ResultStatus
	for _, result := range results {
		if result.TestCaseName == testcase {
			testcaseStatus = result.Status
		}
		utils.CacheSet("testcase_result", result.TestCaseName+"__at__"+pkg, result)
	}
	for _, result := range results {
		if result.TestCaseName != testcase && result.Status == testcaseStatus {
			newFellows = append(newFellows, result.TestCaseName)
		}
	}
	return testcaseStatus == tester.ResultPass, newFellows, nil
}
