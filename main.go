package main

import (
	"fotff/fotff"
	"fotff/pkg"
	"fotff/pkg/dayu200"
	"fotff/test"
	"fotff/test/xdevice"
	"log"
)

func main() {
	var m pkg.Manager = dayu200.Manager{}
	var t test.Tester = xdevice.Tester{}
	var suite = "pts"
	var testPkg string
	for {
		testPkg = m.GetNewerPkg(testPkg)
		dir, err := m.Pkg2Dir(testPkg)
		if err != nil {
			log.Printf("extact package %s to dir err: %v", testPkg, err)
			continue
		}
		if err := m.Flash(dir); err != nil {
			log.Printf("flash package dir %s err: %v", dir, err)
			continue
		}
		results, err := t.DoTestSuite(suite)
		if err != nil {
			log.Printf("do test suite for package %s dir %s err: %v", testPkg, dir, err)
			continue
		}
		fotff.Analysis(m, t, testPkg, results)
	}
}
