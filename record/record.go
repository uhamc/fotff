package record

import (
	"fotff/fotff"
	"fotff/pkg"
	"fotff/test"
	"log"
)

var Records map[string]Record

func HandleResults(m pkg.Manager, pkgName string, results []test.Result) {
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
	handleFailResults(m, pkgName, fails)
}

func handlePassResults(pkgName string, results []test.Result) {
	for _, result := range results {
		Records[result.TestCaseName] = Record{
			LatestResult:     result,
			LatestSuccessPkg: pkgName,
			FailIssueURL:     "",
		}
	}
}

func handleFailResults(m pkg.Manager, pkgName string, results []test.Result) {
	for _, result := range results {
		if Records[result.TestCaseName].FailIssueURL != "" {
			log.Printf("test case %s had failed before and had been welly handled, skip handle it", result.TestCaseName)
			continue
		}
		latestSuccessPkg := Records[result.TestCaseName].LatestSuccessPkg
		Records[result.TestCaseName] = Record{
			LatestResult:     result,
			LatestSuccessPkg: latestSuccessPkg,
			FailIssueURL:     fotff.FindOutTheFirstFail(result.TestCaseName, m, latestSuccessPkg, pkgName),
		}
	}
}
