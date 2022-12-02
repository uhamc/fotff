package xdevice

import "fotff/test"

type Tester struct{}

func (t Tester) DoTestSuite(testSuite string) ([]test.Result, error) {
	panic("TODO")
}

func (t Tester) DoTestCase(testCase string) (test.Result, error) {
	panic("TODO")
}
