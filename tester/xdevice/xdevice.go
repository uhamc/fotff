package xdevice

import (
	"encoding/xml"
	"fmt"
	"fotff/tester"
	"fotff/utils"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

type Tester struct {
	Task          string `key:"task" default:"acts"`
	Config        string `key:"config" default:"./config/user_config.xml"`
	TestCasesPath string `key:"test_cases_path" default:"./testcases"`
	SN            string `key:"sn" default:""`
}

type Report struct {
	XMLName   xml.Name `xml:"testsuites"`
	TestSuite []struct {
		TestCase []struct {
			Name   string `xml:"name,attr"`
			Result string `xml:"result,attr"`
		} `xml:"testcase"`
	} `xml:"testsuite"`
}

func NewTester() tester.Tester {
	ret := &Tester{}
	utils.ParseFromConfigFile("xdevice", ret)
	return ret
}

func (t *Tester) TaskName() string {
	return t.Task
}

func (t *Tester) DoTestTask() (ret []tester.Result, err error) {
	args := []string{"-m", "xdevice", "run", t.Task, "-c", t.Config, "-tcpath", t.TestCasesPath}
	if t.SN != "" {
		args = append(args, "-sn", t.SN)
	}
	if err := utils.Exec("python", args...); err != nil {
		logrus.Errorf("do test suite fail: %v", err)
		return nil, err
	}
	return t.readLatestReport()
}

func (t *Tester) DoTestCase(testCase string) (ret tester.Result, err error) {
	args := []string{"-m", "xdevice", "run", "-l", testCase, "-c", t.Config, "-tcpath", t.TestCasesPath}
	if t.SN != "" {
		args = append(args, "-sn", t.SN)
	}
	if err := utils.Exec("python", args...); err != nil {
		logrus.Errorf("do test case %s fail: %v", testCase, err)
		return ret, err
	}
	r, err := t.readLatestReport()
	if len(r) == 0 {
		return ret, fmt.Errorf("read latest report err, no result found")
	}
	if r[0].TestCaseName != testCase {
		return ret, fmt.Errorf("read latest report err, no matched result found")
	}
	return r[0], nil
}

func (t *Tester) readLatestReport() (ret []tester.Result, err error) {
	data, err := os.ReadFile(filepath.Join("reports", "latest", "summary_report.xml"))
	if err != nil {
		logrus.Errorf("read report xml fail: %v", err)
		return nil, err
	}
	var report Report
	err = xml.Unmarshal(data, &report)
	if err != nil {
		logrus.Errorf("unmarshal report xml fail: %v", err)
		return nil, err
	}
	for _, s := range report.TestSuite {
		for _, c := range s.TestCase {
			var status tester.ResultStatus
			if c.Result == "true" {
				status = tester.ResultPass
			} else {
				status = tester.ResultFail
			}
			ret = append(ret, tester.Result{TestCaseName: c.Name, Status: status})
		}
	}
	return ret, err
}
