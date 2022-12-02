package fotff

import (
	"errors"
	"fmt"
	"fotff/pkg"
	"fotff/test"
	"fotff/vcs"
)

// FindOutTheFirstFail returns the first issue URL that introduce the failure.
func FindOutTheFirstFail(m pkg.Manager, t test.Tester, testCase string, latestSuccessPkg string, failPkg string) (string, error) {
	if latestSuccessPkg == "" {
		return "", fmt.Errorf("can not get the latest success package for %s", testCase)
	}
	from, err := m.GetManifest(latestSuccessPkg)
	if err != nil {
		return "", err
	}
	to, err := m.GetManifest(failPkg)
	if err != nil {
		return "", err
	}
	manifestSteps := vcs.ManifestStepsExpand(from, to)
	return findOutTheFirstFail(m, t, testCase, manifestSteps)
}

func findOutTheFirstFail(m pkg.Manager, t test.Tester, testCase string, steps []vcs.ManifestStep) (string, error) {
	if len(steps) < 2 {
		return "", errors.New("steps are no between a success and a failure")
	}
	if len(steps) == 2 {
		return steps[1].LatestIssueURL, nil
	}
	toTest := len(steps) / 2
	toTestDir, err := m.GenPkgDir(steps[toTest].Manifest)
	if err != nil {
		return "", err
	}
	if err := m.Flash(toTestDir); err != nil {
		return "", err
	}
	result, err := t.DoTestCase(testCase)
	if err != nil {
		return "", err
	}
	if result.Status == test.ResultPass {
		return findOutTheFirstFail(m, t, testCase, steps[toTest:])
	} else {
		return findOutTheFirstFail(m, t, testCase, steps[:toTest+1])
	}
}
