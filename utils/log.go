package utils

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var logOutputFile io.WriteCloser
var logOutputLock sync.RWMutex

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

func LogToStderr() {
	logOutputLock.Lock()
	logrus.SetOutput(os.Stderr)
	if logOutputFile != nil {
		logOutputFile.Close()
		logOutputFile = nil
	}
	logOutputLock.Unlock()
}

func SetLogOutput(pkg string) {
	LogToStderr()
	logOutputLock.Lock()
	defer logOutputLock.Unlock()
	file := filepath.Join("logs", pkg+".log")
	f, err := os.Create(file)
	if err != nil {
		logrus.Errorf("failed to open new log file %s: %v", file, err)
		return
	}
	logrus.Infof("now logs to %s", file)
	logOutputFile = f
	logrus.SetOutput(f)
}

func LogRLock() {
	logOutputLock.RLock()
}

func LogRUnlock() {
	logOutputLock.RUnlock()
}

func GetLogOutput() io.Writer {
	if logOutputFile != nil {
		return logOutputFile
	}
	return os.Stderr
}
