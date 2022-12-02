package dayu200

import (
	"fotff/pkg"
	"fotff/vcs"
	"time"
)

type Manager struct{}

func (m Manager) Flash(dir string) error {
	//TODO implement me
	panic("implement me")
}

func (m Manager) GetManifest(dir string) (vcs.Manifest, error) {
	//TODO implement me
	panic("implement me")
}

func (m Manager) GenPkgDir(manifest vcs.Manifest) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m Manager) GetNewerPkg(pkgName string) string {
	return pkg.GetDirNewerPkg(`%%USERPROFILE%%\Downloads\dayu200`, pkgName)
}

func (m Manager) GetPkgTime(pkgName string) (time.Time, error) {
	//TODO implement me
	panic("implement me")
}

func (m Manager) Pkg2Dir(pkgName string) (string, error) {
	//TODO implement me
	panic("implement me")
}
