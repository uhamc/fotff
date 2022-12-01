package dayu200

import (
	"fotff/pkg"
	"time"
)

type Manager struct{}

func (m Manager) GetNewerPkg(pkgName string) string {
	//TODO implement me
	return pkg.GetDirNewerPkg(`%USERPROFILE%\Downloads\dayu200`, pkgName)
}

func (m Manager) GetManifest(pkgName string) string {
	//TODO implement me
	panic("implement me")
}

func (m Manager) GetTime(pkgName string) time.Time {
	//TODO implement me
	panic("implement me")
}

func (m Manager) Flash(pkgName string) error {
	//TODO implement me
	panic("implement me")
}
