package main

import (
	"fotff/pkg"
	"fotff/pkg/dayu200"
	"fotff/record"
	"fotff/test"
	"fotff/test/xdevice"
)

func main() {
	var m pkg.Manager = dayu200.Manager{}
	var t test.Tester = xdevice.Tester{}
	var suite string = "pts"
	var testPkg string
	for {
		testPkg = m.GetNewerPkg(testPkg)
		results := t.DoTestSuite(suite)
		record.HandleResults(m, testPkg, results)
	}
}
