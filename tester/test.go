package tester

type ResultStatus string

const (
	ResultPass = `pass`
	ResultFail = `fail`
)

type Result struct {
	TestCaseName string
	Status       ResultStatus
}

type Tester interface {
	DoTestTask() ([]Result, error)
	DoTestCase(testCase string) (Result, error)
}

type NewFunc func() Tester
