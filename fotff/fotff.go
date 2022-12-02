package fotff

import (
	"errors"
	"fmt"
	"fotff/pkg"
	"fotff/test"
)

// FindOutTheFirstFail returns the first issue URL that introduce the failure.
func FindOutTheFirstFail(m pkg.Manager, t test.Tester, testCase string, successPkg string, failPkg string) (string, error) {
	if successPkg == "" {
		return "", fmt.Errorf("can not get a success package for %s", testCase)
	}
	steps, err := m.Steps(successPkg, failPkg)
	if err != nil {
		return "", err
	}
	return findOutTheFirstFail(m, t, testCase, steps)
}

func findOutTheFirstFail(m pkg.Manager, t test.Tester, testCase string, steps []string) (string, error) {
	if len(steps) < 2 {
		return "", errors.New("steps are no between a success and a failure")
	}
	if len(steps) == 2 {
		return m.LastIssue(steps[1])
	}
	mid := len(steps) / 2
	if err := m.Flash(steps[mid]); err != nil {
		return "", err
	}
	result, err := t.DoTestCase(testCase)
	if err != nil {
		return "", err
	}
	if result.Status == test.ResultPass {
		return findOutTheFirstFail(m, t, testCase, steps[mid:])
	} else {
		return findOutTheFirstFail(m, t, testCase, steps[:mid+1])
	}
}
