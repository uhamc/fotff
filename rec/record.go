package rec

import (
	"fotff/pkg"
	"fotff/tester"
	"log"
)

var Records = make(map[string]Record)

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
			log.Printf("test case %s had failed before and had been handled, skip handle it", result.TestCaseName)
			continue
		}
		latestSuccessPkg := Records[result.TestCaseName].LatestSuccessPkg
		issueURL, err := FindOutTheFirstFail(m, t, result.TestCaseName, latestSuccessPkg, pkgName)
		if err != nil {
			issueURL = err.Error()
		}
		Records[result.TestCaseName] = Record{
			LatestSuccessPkg: latestSuccessPkg,
			EarliestFailPkg:  pkgName,
			FailIssueURL:     issueURL,
		}
	}
}
