package tester

type ResultStatus string

const (
	ResultPass           = `pass`
	ResultOccasionalFail = `occasional_fail`
	ResultFail           = `fail`
)

type Result struct {
	TestCaseName string
	Status       ResultStatus
}

type Tester interface {
	TaskName() string
	DoTestTask() ([]Result, error)
	DoTestCase(testCase string) (Result, error)
}

type NewFunc func() Tester
