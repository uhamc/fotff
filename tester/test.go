package tester

type Tester interface {
	DoTestSuite() ([]Result, error)
	DoTestCase(testCase string) (Result, error)
}

type NewFunc func() Tester
