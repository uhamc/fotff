package test

type Tester interface {
	DoTestSuite(testSuite string) ([]Result, error)
	DoTestCase(testCase string) (Result, error)
}
