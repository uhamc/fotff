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
		newPkg, err := m.GetNewer()
		if err != nil {
			log.Printf("get newer package err: %v", err)
			continue
		}
		if err := m.Flash(newPkg); err != nil {
			log.Printf("flash package dir %s err: %v", newPkg, err)
			continue
		}
		results, err := t.DoTestSuite(suite)
		if err != nil {
			log.Printf("do test suite for package %s err: %v", newPkg, err)
			continue
		}
		rec.Analysis(m, t, newPkg, results)
	}
}
