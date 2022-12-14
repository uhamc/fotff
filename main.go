package main

import (
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
	m, t := initExecutor()
	data, _ := utils.ReadRuntimeData("last_handled.rec")
	var curPkg = string(data)
	for {
		utils.LogToStdout()
		if err := utils.WriteRuntimeData("last_handled.rec", []byte(curPkg)); err != nil {
			logrus.Errorf("failed to write last_handled.rec: %v", err)
		}
		logrus.Info("waiting for a newer package...")
		var err error
		curPkg, err = m.GetNewer(curPkg)
		if err != nil {
			logrus.Infof("get newer package err: %v", err)
			continue
		}
		utils.SetLogOutput(filepath.Base(curPkg))
		logrus.Infof("now flash %s...", curPkg)
		if err := m.Flash(curPkg); err != nil {
			logrus.Errorf("flash package dir %s err: %v", curPkg, err)
			continue
		}
		logrus.Info("now do test suite...")
		results, err := t.DoTestTask()
		if err != nil {
			logrus.Errorf("do test suite for package %s err: %v", curPkg, err)
			continue
		}
		logrus.Infof("now analysis test results...")
		rec.Analysis(m, t, curPkg, results)
		rec.Report(curPkg, t.TaskName())
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
