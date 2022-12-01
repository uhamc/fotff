package pkg

import (
	"log"
	"os"
	"sort"
	"time"
)

type Manager interface {
	// GetNewerPkg blocks the process until a newer package is found, then returns the newest one.
	GetNewerPkg(pkgName string) string
	// GetManifest returns the manifest info of given package.
	GetManifest(pkgName string) string
	// GetTime returns the creation time of given package.
	GetTime(pkgName string) time.Time
	// Flash download given package to the device.
	Flash(pkgName string) error
}

func GetDirNewerPkg(dir string, pkgName string) string {
	for {
		files, _ := os.ReadDir(dir)
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
