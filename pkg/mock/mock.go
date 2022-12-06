package mock

import (
	"fotff/pkg/dayu200"
	"github.com/sirupsen/logrus"
)

type Manager struct {
	dayu200.Manager
}

func (m *Manager) Flash(pkg string) error {
	m.Manager.Flash(pkg)
	logrus.Warn("mock implementation ignores any error and returns OK unconditionally")
	return nil
}
