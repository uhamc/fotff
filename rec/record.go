package rec

import (
	"encoding/json"
	"fotff/pkg"
	"fotff/tester"
	"github.com/sirupsen/logrus"
	"os"
)

var Records = make(map[string]Record)

func init() {
	if _, err := os.Stat("records.json"); err != nil && os.IsNotExist(err) {
		return
	}
	data, err := os.ReadFile("records.json")
	if err != nil {
		logrus.Errorf("read records.json err: %v", err)
	}
	if err := json.Unmarshal(data, &Records); err != nil {
		logrus.Errorf("unmarshal records err: %v", err)
	}
}

func Save() {
	data, err := json.MarshalIndent(Records, "", "\t")
	if err != nil {
		logrus.Errorf("marshal records err: %v", err)
	}
	if err := os.WriteFile("records.json", data, 0640); err != nil {
		logrus.Errorf("save records err: %v", err)
	}
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
		logrus.Infof("recording %s as a success, the lastest success package is %s", result.TestCaseName, pkgName)
		Records[result.TestCaseName] = Record{
			LatestSuccessPkg: pkgName,
			EarliestFailPkg:  "",
			FailIssueURL:     "",
		}
	}
}

func handleFailResults(m pkg.Manager, t tester.Tester, pkgName string, results []tester.Result) {
	for _, result := range results {
		if Records[result.TestCaseName].EarliestFailPkg != "" {
			logrus.Warnf("test case %s had failed before and had been handled, skip handle it", result.TestCaseName)
			continue
		}
		latestSuccessPkg := Records[result.TestCaseName].LatestSuccessPkg
		var issueURL string
		if latestSuccessPkg != "" {
			var err error
			logrus.Warnf("%s failed, the lastest success package is %s, earliest fail package is %s, now finding out the first fail...", result.TestCaseName, latestSuccessPkg, pkgName)
			issueURL, err = FindOutTheFirstFail(m, t, result.TestCaseName, latestSuccessPkg, pkgName)
			if err != nil {
				logrus.Errorf("failed to find out the first fail issue, err: %v", err)
				issueURL = err.Error()
			}
		}
		logrus.Warnf("recording %s as a failure, the lastest success package is %s, the earliest fail package is %s, fail issue URL is %s", result.TestCaseName, latestSuccessPkg, pkgName, issueURL)
		Records[result.TestCaseName] = Record{
			LatestSuccessPkg: latestSuccessPkg,
			EarliestFailPkg:  pkgName,
			FailIssueURL:     issueURL,
		}
	}
}
