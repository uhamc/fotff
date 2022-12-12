package rec

type Record struct {
	UpdateTime       string `col:"update time"`
	Status           string `col:"status"`
	LatestSuccessPkg string `col:"last success package"`
	EarliestFailPkg  string `col:"earliest fail package"`
	FailIssueURL     string `col:"fail issue url"`
}
