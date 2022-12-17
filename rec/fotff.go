package rec

import (
	"context"
	"errors"
	"fmt"
	"fotff/pkg"
	"fotff/res"
	"fotff/tester"
	"github.com/sirupsen/logrus"
	"math"
	"sync"
)

type cancelCtx struct {
	ctx context.Context
	fn  context.CancelFunc
}

// FindOutTheFirstFail returns the first issue URL that introduce the failure.
func FindOutTheFirstFail(m pkg.Manager, t tester.Tester, testCase string, successPkg string, failPkg string) (string, error) {
	if successPkg == "" {
		return "", fmt.Errorf("can not get a success package for %s", testCase)
	}
	steps, err := m.Steps(successPkg, failPkg)
	if err != nil {
		return "", err
	}
	return findOutTheFirstFail(m, t, testCase, steps)
}

func findOutTheFirstFail(m pkg.Manager, t tester.Tester, testcase string, steps []string) (string, error) {
	if len(steps) == 0 {
		return "", errors.New("steps are no between (success, failure]")
	}
	logrus.Infof("now use %d-section search to find out the first fault, the length of range is %d, between [%s, %s]", res.Num(), len(steps), steps[0], steps[len(steps)-1])
	if len(steps) == 1 {
		return m.LastIssue(steps[0])
	}
	success, fail := -1, len(steps)-1
	var lock sync.Mutex
	gapLen := float64(len(steps)-1) / float64(res.Num()+1)
	if gapLen < 1 {
		gapLen = 1
	}
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
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 1; i <= res.Num(); i++ {
		index := len(steps) - 1 - int(math.Round(float64(i)*gapLen)) // index from the tail to avoid testing the last one, on which can not narrow ranges
		if index < 0 {
			break
		}
		ctx, fn := context.WithCancel(context.WithValue(context.TODO(), "index", index))
		contexts = append(contexts, cancelCtx{ctx: ctx, fn: fn})
		wg.Add(1)
		go func(index int, ctx context.Context) {
			defer wg.Done()
			<-start
			pass, err := flashAndTest(m, t, steps[index], testcase, ctx)
			if err != nil {
				if !errors.Is(err, context.Canceled) {
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
	return findOutTheFirstFail(m, t, testcase, steps[success+1:fail+1])
}

func flashAndTest(m pkg.Manager, t tester.Tester, pkg string, testcase string, ctx context.Context) (bool, error) {
	device := res.GetDevice()
	defer res.ReleaseDevice(device)
	if err := m.Flash(device, pkg, ctx); err != nil {
		return false, err
	}
	result, err := t.DoTestCase(device, testcase, ctx)
	if err != nil {
		return false, err
	}
	return result.Status == tester.ResultPass, nil
}
