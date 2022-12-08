package xdevice

import (
	"encoding/json"
	"fotff/tester"
	"fotff/utils"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

type Tester struct {
	RunTestSuitePy string `key:"run_test_suite_py" default:"./runallcase.py"`
	RunTestCasePy  string `key:"run_test_case_py" default:"./runonecase.py"`
	ResultJson     string `key:"result_json" default:"./result.json"`
}

func NewTester() tester.Tester {
	ret := &Tester{}
	utils.ParseFromConfigFile("xdevice", ret)
	return ret
}

func (t *Tester) DoTestSuite() (ret []tester.Result, err error) {
	out, err := exec.Command("python", t.RunTestCasePy).CombinedOutput()
	if err != nil {
		logrus.Errorf("%s", string(out))
		logrus.Errorf("do test suite fail: %v", err)
		return nil, err
	}
	data, err := os.ReadFile(t.ResultJson)
	if err != nil {
		logrus.Errorf("read result json err: %v", err)
	}
	err = json.Unmarshal(data, &ret)
	return ret, err
}

func (t *Tester) DoTestCase(testCase string) (ret tester.Result, err error) {
	out, err := exec.Command("python", t.RunTestCasePy, testCase).CombinedOutput()
	if err != nil {
		logrus.Errorf("%s", string(out))
		logrus.Errorf("do test case %s fail: %v", testCase, err)
		return tester.Result{}, err
	}
	data, err := os.ReadFile(t.ResultJson)
	if err != nil {
		logrus.Errorf("read result json err: %v", err)
	}
	err = json.Unmarshal(data, &ret)
	return ret, err
}
