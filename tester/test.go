package tester

import "context"

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
	DoTestTask(device string, ctx context.Context) ([]Result, error)
	DoTestCase(device string, testCase string, ctx context.Context) (Result, error)
}

type NewFunc func() Tester
