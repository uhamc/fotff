package mock

import (
	"fotff/pkg"
	"fotff/pkg/dayu200"
	"github.com/sirupsen/logrus"
)

type Manager struct {
	dayu200.Manager
}

func NewManager() pkg.Manager {
	return &Manager{
		Manager: dayu200.Manager{
			PkgDir:    `C:\dayu200`,
			Workspace: `C:\dayu200_workspace`,
			Branch:    "master",
			BuildServerConfig: dayu200.BuildServerConfig{
				Addr:           "172.0.0.1:22",
				User:           "sample",
				Passwd:         "samplePasswd",
				BuildWorkSpace: "/home/sample/fotff/build_workspace",
			},
		},
	}
}

func (m *Manager) Flash(pkg string) error {
	m.Manager.Flash(pkg)
	logrus.Warn("mock implementation ignores any error and returns OK unconditionally")
	return nil
}
