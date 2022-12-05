package main

import (
	"fmt"
	"fotff/pkg"
	"fotff/pkg/dayu200"
	"fotff/pkg/mock"
	"fotff/rec"
	"fotff/tester"
	testermock "fotff/tester/mock"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	initLogrus()
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
			logrus.Infof("get newer package err: %v", err)
			continue
		}
		logrus.Infof("now flash %s...", newPkg)
		if err := m.Flash(newPkg); err != nil {
			logrus.Errorf("flash package dir %s err: %v", newPkg, err)
			continue
		}
		logrus.Infof("now do test suite %s...", suite)
		results, err := t.DoTestSuite(suite)
		if err != nil {
			logrus.Errorf("do test suite for package %s err: %v", newPkg, err)
			continue
		}
		logrus.Infof("now analysis test results of %s...", suite)
		rec.Analysis(m, t, newPkg, results)
	}
}

func initLogrus() {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			funcName := strings.Split(f.Function, ".")
			fn := funcName[len(funcName)-1]
			_, filename := filepath.Split(f.File)
			return fmt.Sprintf("%s()", fn), fmt.Sprintf("%s:%d", filename, f.Line)
		},
	})
}
