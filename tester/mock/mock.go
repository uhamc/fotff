package mock

import (
	"fotff/rec"
	"fotff/tester"
	"fotff/tester/xdevice"
	"github.com/sirupsen/logrus"
)

type Tester struct {
	xdevice.Tester
}

func NewTester() tester.Tester {
	return &Tester{}
}

func init() {
	rec.Records["TEST_001"] = rec.Record{
		LatestSuccessPkg: `C:\dayu200_workspace\version-Daily_Version-dayu200-20221201_072124-dayu200`,
		EarliestFailPkg:  ``,
		FailIssueURL:     "",
	}
	rec.Records["TEST_002"] = rec.Record{
		LatestSuccessPkg: `C:\dayu200_workspace\version-Daily_Version-dayu200-20221201_072124-dayu200`,
		EarliestFailPkg:  ``,
		FailIssueURL:     "",
	}
}

func (t Tester) DoTestSuite() ([]tester.Result, error) {
	logrus.Infof("TEST_001 pass")
	logrus.Warnf("TEST_002 fail")
	return []tester.Result{
		{TestCaseName: "TEST_001", Status: tester.ResultPass},
		{TestCaseName: "TEST_002", Status: tester.ResultFail},
	}, nil
}

func (t Tester) DoTestCase(testCase string) (tester.Result, error) {
	switch testCase {
	case "TEST_001":
		logrus.Infof("TEST_001 pass")
		return tester.Result{TestCaseName: "TEST_001", Status: tester.ResultPass}, nil
	case "TEST_002":
		logrus.Warnf("TEST_002 fail")
		return tester.Result{TestCaseName: "TEST_002", Status: tester.ResultFail}, nil
	default:
		panic("not defined")
	}
}
