package xdevice

import "fotff/tester"

type Tester struct{}

func (t Tester) DoTestSuite(testSuite string) ([]tester.Result, error) {
	panic("TODO")
}

func (t Tester) DoTestCase(testCase string) (tester.Result, error) {
	panic("TODO")
}