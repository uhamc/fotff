package rec

import (
	"encoding/json"
	"fotff/pkg"
	"fotff/tester"
	"fotff/utils"
	"github.com/sirupsen/logrus"
	"time"
)

var Records = make(map[string]Record)

func init() {
	data, err := utils.ReadRuntimeData("records.json")
	if err != nil {
		return
	}
	if err := json.Unmarshal(data, &Records); err != nil {
		logrus.Errorf("unmarshal records err: %v", err)
	}
}

func Save() {
	data, err := json.MarshalIndent(Records, "", "\t")
	if err != nil {
		logrus.Errorf("marshal records err: %v", err)
		return
	}
	if err := utils.WriteRuntimeData("records.json", data); err != nil {
		logrus.Errorf("save records err: %v", err)
		return
	}
	logrus.Infof("save records successfully")
}

func Analysis(m pkg.Manager, t tester.Tester, pkgName string, results []tester.Result) {
	var passes, fails []tester.Result
	for _, result := range results {
		switch result.Status {
		case tester.ResultPass:
			passes = append(passes, result)
		case tester.ResultFail:
			fails = append(fails, result)
		}
	}
	handlePassResults(pkgName, passes)
	handleFailResults(m, t, pkgName, fails)
}

func handlePassResults(pkgName string, results []tester.Result) {
	for _, result := range results {
		logrus.Infof("recording [%s] as a success, the lastest success package is [%s]", result.TestCaseName, pkgName)
		Records[result.TestCaseName] = Record{
			UpdateTime:       time.Now().Format("2006-01-02 15:04:05"),
			Status:           tester.ResultPass,
			LatestSuccessPkg: pkgName,
			EarliestFailPkg:  "",
			FailIssueURL:     "",
		}
	}
}

func handleFailResults(m pkg.Manager, t tester.Tester, pkgName string, results []tester.Result) {
loop:
	for _, result := range results {
		if record, ok := Records[result.TestCaseName]; ok && record.Status != tester.ResultPass {
			logrus.Warnf("test case %s had failed before, skip handle it", result.TestCaseName)
			continue
		}
		latestSuccessPkg := Records[result.TestCaseName].LatestSuccessPkg
		for i := 0; i < 2; i++ {
			r, err := t.DoTestCase(result.TestCaseName)
			if err != nil {
				logrus.Errorf("failed to do test case %s: %v", result.TestCaseName, err)
				continue
			}
			if r.Status == tester.ResultPass {
				Records[result.TestCaseName] = Record{
					UpdateTime:       time.Now().Format("2006-01-02 15:04:05"),
					Status:           tester.ResultOccasionalFail,
					LatestSuccessPkg: latestSuccessPkg,
					EarliestFailPkg:  pkgName,
					FailIssueURL:     "seems to be an occasional issue, skip analysing",
				}
				continue loop
			}
		}
		var issueURL string
		if latestSuccessPkg != "" {
			var err error
			logrus.Warnf("%s failed, the lastest success package is [%s], earliest fail package is [%s], now finding out the first fail...", result.TestCaseName, latestSuccessPkg, pkgName)
			issueURL, err = FindOutTheFirstFail(m, t, result.TestCaseName, latestSuccessPkg, pkgName)
			if err != nil {
				logrus.Errorf("failed to find out the first fail issue, err: %v", err)
				issueURL = err.Error()
			}
		} else {
			issueURL = "no previous success found, can not analysis"
		}
		logrus.Warnf("recording %s as a failure, the lastest success package is [%s], the earliest fail package is [%s], fail issue URL is [%s]", result.TestCaseName, latestSuccessPkg, pkgName, issueURL)
		Records[result.TestCaseName] = Record{
			UpdateTime:       time.Now().Format("2006-01-02 15:04:05"),
			Status:           tester.ResultFail,
			LatestSuccessPkg: latestSuccessPkg,
			EarliestFailPkg:  pkgName,
			FailIssueURL:     issueURL,
		}
	}
}
