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
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"io"
	"os/exec"
	"time"
)

func ExecContext(ctx context.Context, name string, args ...string) error {
	if err := execContext(ctx, name, args...); err != nil {
		if errors.Is(err, context.Canceled) {
			return err
		}
		logrus.Errorf("exec failed: %v, try again...", err)
		return execContext(ctx, name, args...)
	}
	return nil
}

func execContext(ctx context.Context, name string, args ...string) error {
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

func ExecCombinedOutputContext(ctx context.Context, name string, args ...string) ([]byte, error) {
	out, err := execCombinedOutputContext(ctx, name, args...)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return out, err
		}
		logrus.Errorf("exec failed: %v, try again...", err)
		return execCombinedOutputContext(ctx, name, args...)
	}
	return out, nil
}

func execCombinedOutputContext(ctx context.Context, name string, args ...string) ([]byte, error) {
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
