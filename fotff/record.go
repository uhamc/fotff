package fotff

import (
	"fotff/pkg"
	"fotff/test"
	"log"
)

var Records map[string]Record

func Analysis(m pkg.Manager, t test.Tester, pkgName string, results []test.Result) {
	var passes, fails []test.Result
	for _, result := range results {
		switch result.Status {
		case test.ResultPass:
			passes = append(passes, result)
		case test.ResultFail:
			fails = append(fails, result)
		}
	}
	handlePassResults(pkgName, passes)
	handleFailResults(m, t, pkgName, fails)
}

func handlePassResults(pkgName string, results []test.Result) {
	for _, result := range results {
		Records[result.TestCaseName] = Record{
			LatestSuccessPkg: pkgName,
			EarliestFailPkg:  "",
			FailIssueURL:     "",
		}
	}
}

func handleFailResults(m pkg.Manager, t test.Tester, pkgName string, results []test.Result) {
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
