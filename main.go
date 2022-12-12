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
	"fotff/utils"
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
	data, _ := utils.ReadRuntimeData("last_handled.rec")
	var curPkg = string(data)
	for {
		if err := utils.WriteRuntimeData("last_handled.rec", []byte(curPkg)); err != nil {
			logrus.Errorf("failed to write last_handled.rec: %v", err)
		}
		newPkg, err := m.GetNewer(curPkg)
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
		results, err := t.DoTestTask()
		if err != nil {
			logrus.Errorf("do test suite for package %s err: %v", newPkg, err)
			continue
		}
		logrus.Infof("now analysis test results...")
		rec.Analysis(m, t, newPkg, results)
		rec.Save()
		curPkg = newPkg
	}
}

func initExecutor() (pkg.Manager, tester.Tester) {
	//TODO load from config file
	var conf = struct {
		PkgManager string `key:"pkg_manager" default:"mock"`
		Tester     string `key:"tester" default:"mock"`
	}{}
	utils.ParseFromConfigFile("", &conf)
	newPkgMgrFunc, ok := newPkgMgrFuncs[conf.PkgManager]
	if !ok {
		logrus.Panicf("no package manager found for %s", conf.PkgManager)
	}
	newTesterFunc, ok := newTesterFuncs[conf.Tester]
	if !ok {
		logrus.Panicf("no tester found for %s", conf.Tester)
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
