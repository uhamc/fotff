package test

type Tester interface {
	DoTestSuite(testSuite string) []Result
	DoTestCase(testCase string) Result
}
