package mock

import "fotff/pkg/dayu200"

type Manager struct {
	dayu200.Manager
}

func (m *Manager) Flash(pkg string) error {
	return nil
}
