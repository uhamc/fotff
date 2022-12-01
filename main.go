package main

import (
	"fotff/pkg"
	"fotff/pkg/dayu200"
	"fotff/record"
	"fotff/test/xdevice"
)

func main() {
	var m pkg.Manager = dayu200.Manager{}
	var testPkg string
	for {
		testPkg = m.GetNewerPkg(testPkg)
		results := xdevice.DoFullTest(testPkg)
		record.HandleResults(m, testPkg, results)
	}
}
