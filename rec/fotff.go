package rec

import (
	"errors"
	"fmt"
	"fotff/pkg"
	"fotff/tester"
	"github.com/sirupsen/logrus"
)

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

func findOutTheFirstFail(m pkg.Manager, t tester.Tester, testCase string, steps []string) (string, error) {
	if len(steps) == 0 {
		return "", errors.New("steps are no between (success, failure]")
	}
	logrus.Infof("now bin-search the first failed step, between [%s, %s]", steps[0], steps[len(steps)-1])
	if len(steps) == 1 {
		return m.LastIssue(steps[0])
	}
	mid := (len(steps) - 1) / 2
	if err := m.Flash(steps[mid]); err != nil {
		return "", err
	}
	result, err := t.DoTestCase(testCase)
	if err != nil {
		return "", err
	}
	if result.Status == tester.ResultPass {
		return findOutTheFirstFail(m, t, testCase, steps[mid+1:])
	} else {
		return findOutTheFirstFail(m, t, testCase, steps[:mid+1])
	}
}
