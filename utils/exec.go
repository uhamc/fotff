package utils

import (
	"github.com/sirupsen/logrus"
	"io"
	"os/exec"
)

func Exec(name string, args ...string) error {
	LogRLock()
	defer LogRUnlock()
	cmdStr := append([]string{name}, args...)
	logrus.Infof("cmd: %s", cmdStr)
	cmd := exec.Command(name, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	go io.Copy(GetLogOutput(), stdout)
	go io.Copy(GetLogOutput(), stderr)
	return cmd.Wait()
}

func ExecCombinedOutput(name string, args ...string) ([]byte, error) {
	cmdStr := append([]string{name}, args...)
	logrus.Infof("cmd: %s", cmdStr)
	out, err := exec.Command(name, args...).CombinedOutput()
	logrus.Infof("out: %s", string(out))
	return out, err
}
