package mock

import (
	"fmt"
	"fotff/pkg"
	"github.com/sirupsen/logrus"
	"time"
)

type Manager struct {
	pkgCount int
}

func NewManager() pkg.Manager {
	return &Manager{}
}

func (m *Manager) LastIssue(pkg string) (string, error) {
	ret := fmt.Sprintf("https://testserver.com/issues/%s", pkg)
	logrus.Infof("LastIssue: mock implementation returns %s", ret)
	return ret, nil
}

func (m *Manager) Steps(from, to string) ([]string, error) {
	var ret = []string{"step1", "step2", "step3"}
	for i := range ret {
		ret[i] = fmt.Sprintf("%s-%s-%s", from, to, ret[i])
	}
	logrus.Infof("Steps: mock implementation returns %v", ret)
	return ret, nil
}

func (m *Manager) GetNewer() (string, error) {
	ret := fmt.Sprintf("pkg%d", m.pkgCount)
	time.Sleep(time.Duration(m.pkgCount) * time.Second)
	m.pkgCount++
	logrus.Infof("GetNewer: mock implementation returns %s", ret)
	return ret, nil
}

func (m *Manager) Flash(pkg string) error {
	time.Sleep(time.Second)
	logrus.Infof("Flash: mock implementation returns OK unconditionally")
	return nil
}
