/*
 * Copyright (c) 2022 Huawei Device Co., Ltd.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
