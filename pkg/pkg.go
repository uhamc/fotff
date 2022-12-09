package pkg

import (
	"github.com/sirupsen/logrus"
	"os"
	"sort"
	"time"
)

type NewFunc func() Manager

type Manager interface {
	// Flash download given package dir to the device.
	Flash(pkg string) error
	// LastIssue returns the last issue URL related to the package.
	LastIssue(pkg string) (string, error)
	// Steps generates every intermediate package and returns the list sequentially.
	Steps(from, to string) ([]string, error)
	// GetNewer blocks the process until a newer package is found, then returns the newest one.
	GetNewer() (string, error)
}

func GetNewerFileFromDir(dir string, fileName string, less func(files []os.DirEntry, i, j int) bool) string {
	for {
		files, err := os.ReadDir(dir)
		if err != nil {
			logrus.Errorf("read dir %s err: %s", dir, err)
			time.Sleep(time.Second)
			continue
		}
		sort.Slice(files, func(i, j int) bool {
			return less(files, i, j)
		})
		if len(files) != 0 {
			f := files[len(files)-1]
			if f.Name() != fileName {
				logrus.Infof("new package found, name: %s", f.Name())
				return f.Name()
			}
		}
		time.Sleep(time.Second)
	}
}
