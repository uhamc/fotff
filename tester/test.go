package tester

type Tester interface {
	DoTestTask() ([]Result, error)
	DoTestCase(testCase string) (Result, error)
}

type NewFunc func() Tester
