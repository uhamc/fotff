package main

import (
	"fotff/pkg"
	"fotff/pkg/dayu200"
	"fotff/pkg/mock"
	"fotff/rec"
	"fotff/tester"
	testermock "fotff/tester/mock"
	"log"
)

func main() {
	var m pkg.Manager = &mock.Manager{
		Manager: dayu200.Manager{
			PkgDir:    `C:\dayu200`,
			Workspace: `C:\dayu200_workspace`,
		},
	}
	var t tester.Tester = testermock.Tester{}
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
		rec.Analysis(m, t, pkg, results)
	}
}
