package main

import (
	"fmt"
	"fotff/pkg"
	"fotff/pkg/dayu200"
	"fotff/pkg/mock"
	"fotff/rec"
	"fotff/tester"
	testermock "fotff/tester/mock"
	"fotff/tester/xdevice"
	"github.com/Unknwon/goconfig"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"runtime"
	"strings"
)

var newPkgMgrFuncs = map[string]pkg.NewFunc{
	"mock":    mock.NewManager,
	"dayu200": dayu200.NewManager,
}

var newTesterFuncs = map[string]tester.NewFunc{
	"mock":    testermock.NewTester,
	"xdevice": xdevice.NewTester,
}

func main() {
	initLogrus()
	m, t := initExecutor()
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
		logrus.Info("now do test suite...")
		results, err := t.DoTestSuite()
		if err != nil {
			logrus.Errorf("do test suite for package %s err: %v", newPkg, err)
			continue
		}
		logrus.Infof("now analysis test results...")
		rec.Analysis(m, t, newPkg, results)
		rec.Save()
	}
}

func initExecutor() (pkg.Manager, tester.Tester) {
	//TODO load from config file
	pkgManagerType := "mock"
	testerType := "mock"
	conf, err := goconfig.LoadConfigFile("config.ini")
	if err != nil {
		logrus.Errorf("load config file err, use 'mock' by default: %v", err)
	} else {
		if v, err := conf.GetValue("", "pkg_manager"); err != nil {
			logrus.Errorf("get pkg_manager err, use 'mock' by default: %v", err)
		} else {
			pkgManagerType = v
		}
		if v, err := conf.GetValue("", "tester"); err != nil {
			logrus.Errorf("get tester err, use 'mock' by default: %v", err)
		} else {
			testerType = v
		}
	}
	newPkgMgrFunc, ok := newPkgMgrFuncs[pkgManagerType]
	if !ok {
		logrus.Panicf("no package manager found for %s", pkgManagerType)
	}
	newTesterFunc, ok := newTesterFuncs[testerType]
	if !ok {
		logrus.Panicf("no package manager found for %s", pkgManagerType)
	}
	return newPkgMgrFunc(), newTesterFunc()
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
