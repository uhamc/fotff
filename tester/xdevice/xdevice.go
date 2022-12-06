package xdevice

import "fotff/tester"

type Tester struct{}

func NewTester() tester.Tester {
	return &Tester{}
}

func (t *Tester) DoTestSuite() ([]tester.Result, error) {
	panic("TODO")
}

func (t *Tester) DoTestCase(testCase string) (tester.Result, error) {
	panic("TODO")
}
