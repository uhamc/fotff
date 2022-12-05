package tester

type ResultStatus string

const (
	ResultPass = `Pass`
	ResultFail = `Fail`
)

type Result struct {
	TestCaseName string
	Status       ResultStatus
}
