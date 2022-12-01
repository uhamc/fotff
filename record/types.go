package record

import "fotff/test"

type Record struct {
	LatestResult     test.Result
	LatestSuccessPkg string
	FailIssueURL     string
}
