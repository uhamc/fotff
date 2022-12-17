package main

import (
	"context"
	"fotff/pkg"
	"fotff/pkg/dayu200"
	"fotff/pkg/mock"
	"fotff/rec"
	"fotff/res"
	"fotff/tester"
	testermock "fotff/tester/mock"
	"fotff/tester/xdevice"
	"fotff/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
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

var rootCmd *cobra.Command

func init() {
	m, t := initExecutor()
	rootCmd = &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			loop(m, t)
		},
	}
	runCmd := initRunCmd(m, t)
	flashCmd := initFlashCmd(m)
	rootCmd.AddCommand(runCmd, flashCmd)
}

func initRunCmd(m pkg.Manager, t tester.Tester) *cobra.Command {
	var success, fail, testcase string
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "bin-search in (success, fail] by do given testcase to find out the fist fail, and print the corresponding issue",
		RunE: func(cmd *cobra.Command, args []string) error {
			if success == "" || fail == "" || testcase == "" {

			}
			return fotff(m, t, success, fail, testcase)
		},
	}
	runCmd.PersistentFlags().StringVarP(&success, "success", "s", "", "success package directory")
	runCmd.PersistentFlags().StringVarP(&fail, "fail", "f", "", "fail package directory")
	runCmd.PersistentFlags().StringVarP(&testcase, "testcase", "t", "", "testcase name")
	runCmd.MarkPersistentFlagRequired("success")
	runCmd.MarkPersistentFlagRequired("fail")
	runCmd.MarkPersistentFlagRequired("testcase")
	return runCmd
}

func initFlashCmd(m pkg.Manager) *cobra.Command {
	var flashPkg, device string
	flashCmd := &cobra.Command{
		Use:   "flash",
		Short: "flash the given package",
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.Flash(device, flashPkg, context.TODO())
		},
	}
	flashCmd.PersistentFlags().StringVarP(&flashPkg, "package", "p", "", "package directory")
	flashCmd.PersistentFlags().StringVarP(&flashPkg, "device", "d", "", "device sn")
	flashCmd.MarkPersistentFlagRequired("package")
	return flashCmd
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Errorf("failed to execute: %v", err)
		os.Exit(1)
	}
}

func loop(m pkg.Manager, t tester.Tester) {
	data, _ := utils.ReadRuntimeData("last_handled.rec")
	var curPkg = string(data)
	for {
		utils.LogToStderr()
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
		device := res.GetDevice()
		if err := m.Flash(device, curPkg, context.TODO()); err != nil {
			logrus.Errorf("flash package dir %s err: %v", curPkg, err)
			res.ReleaseDevice(device)
			continue
		}
		logrus.Info("now do test suite...")
		results, err := t.DoTestTask(device, context.TODO())
		if err != nil {
			logrus.Errorf("do test suite for package %s err: %v", curPkg, err)
			continue
		}
		logrus.Infof("now analysis test results...")
		toFotff := rec.HandleResults(t, device, curPkg, results)
		res.ReleaseDevice(device)
		rec.Analysis(m, t, curPkg, toFotff)
		rec.Save()
		rec.Report(curPkg, t.TaskName())
	}
}

func fotff(m pkg.Manager, t tester.Tester, success, fail, testcase string) error {
	issueURL, err := rec.FindOutTheFirstFail(m, t, testcase, success, fail)
	if err != nil {
		logrus.Errorf("failed to find out the first fail: %v", err)
		return err
	}
	logrus.Infof("the first fail found: %v", issueURL)
	return nil
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
