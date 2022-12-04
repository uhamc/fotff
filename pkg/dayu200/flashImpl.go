package dayu200

import (
	"os"
	"path/filepath"
)

func build(pkg string) error {
	if _, err := os.Stat(filepath.Join(pkg, "__to_be_built__")); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	//TODO build package with generated manifest_tag.xml
	panic("implement me")
}
