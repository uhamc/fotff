package mock

import (
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

func (t *Tester) DoTestTask() ([]tester.Result, error) {
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

func (t *Tester) DoTestCase(testCase string) (tester.Result, error) {
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
