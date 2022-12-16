package utils

import (
	"context"
	"github.com/sirupsen/logrus"
	"io"
	"os/exec"
	"time"
)

func Exec(name string, args ...string) error {
	return ExecContext(context.TODO(), name, args...)
}

func ExecContext(ctx context.Context, name string, args ...string) error {
	LogRLock()
	defer LogRUnlock()
	cmdStr := append([]string{name}, args...)
	logrus.Infof("cmd: %s", cmdStr)
	cmd := exec.CommandContext(ctx, name, args...)
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
	return ExecCombinedOutputContext(context.TODO(), name, args...)
}

func ExecCombinedOutputContext(ctx context.Context, name string, args ...string) ([]byte, error) {
	cmdStr := append([]string{name}, args...)
	logrus.Infof("cmd: %s", cmdStr)
	out, err := exec.CommandContext(ctx, name, args...).CombinedOutput()
	logrus.Infof("out: %s", string(out))
	return out, err
}

func SleepContext(duration time.Duration, ctx context.Context) {
	select {
	case <-time.NewTimer(duration).C:
	case <-ctx.Done():
	}
}
