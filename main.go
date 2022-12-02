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
	var m pkg.Manager = &dayu200.Manager{
		PkgDir:    `C:\dayu200`,
		Workspace: `C:\dayu200_workspace`,
	}
	var t test.Tester = xdevice.Tester{}
	var suite = "pts"
	for {
		pkg, err := m.GetNewer()
		if err != nil {
			log.Printf("get newer package err: %v", err)
			continue
		}
		if err := m.Flash(pkg); err != nil {
			log.Printf("flash package dir %s err: %v", pkg, err)
			continue
		}
		results, err := t.DoTestSuite(suite)
		if err != nil {
			log.Printf("do test suite for package %s err: %v", pkg, err)
			continue
		}
		fotff.Analysis(m, t, pkg, results)
	}
}
