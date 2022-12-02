package pkg

import (
	"fotff/vcs"
	"log"
	"os"
	"sort"
	"time"
)

type Manager interface {
	// Flash download given package dir to the device.
	Flash(dir string) error
	// GetManifest returns the manifest info of given package dir.
	GetManifest(dir string) (vcs.Manifest, error)
	// GenPkgDir build a package dir with given manifest.
	GenPkgDir(m vcs.Manifest) (string, error)
	// GetNewerPkg blocks the process until a newer package is found, then returns the newest one.
	GetNewerPkg(pkgName string) string
	// GetPkgTime returns the creation time of given package.
	GetPkgTime(pkgName string) (time.Time, error)
	// Pkg2Dir extracted the package and the path of dir where the package extracted to.
	Pkg2Dir(pkgName string) (string, error)
}

func GetDirNewerPkg(dir string, pkgName string) string {
	for {
		files, err := os.ReadDir(dir)
		if err != nil {
			log.Printf("read dir %s err: %s", dir, err)
			time.Sleep(time.Second)
			continue
		}
		sort.Slice(files, func(i, j int) bool {
			return files[i].Name() < files[j].Name()
		})
		if len(files) != 0 {
			f := files[len(files)-1]
			if f.Name() != pkgName {
				log.Printf("new package found, name: %s", f.Name())
				return f.Name()
			}
		}
		time.Sleep(time.Second)
	}
}
