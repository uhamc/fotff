//go:build linux

package dayu200

import (
	"os"
	"path/filepath"
)

func build(pkg string) error {
	if _, err := os.Stat(filepath.Join(pkg, "__built__")); err == nil {
		return nil
	}
	//TODO build package with generated manifest_tag.xml
	panic("implement me")
}
