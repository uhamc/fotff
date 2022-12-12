package utils

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var logOutput io.WriteCloser

func init() {
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
	os.MkdirAll("logs", 0750)
}

func LogToStdout() {
	logrus.SetOutput(os.Stdout)
	if logOutput != nil {
		logOutput.Close()
	}
}

func SetLogOutput(pkg string) {
	LogToStdout()
	file := filepath.Join("logs", pkg+".log")
	f, err := os.Create(file)
	if err != nil {
		logrus.Errorf("failed to open new log file %s: %v", file, err)
		return
	}
	logrus.Infof("now logs to %s", file)
	logOutput = f
	logrus.SetOutput(f)
}
