package vcs

import (
	"encoding/xml"
	"os"
	"time"
)

type Manifest struct {
	Default  Default   `xml:"default"`
	Projects []Project `xml:"project"`
}

type Default struct {
}

type Project struct {
	Name     string `xml:"name,attr"`
	Path     string `xml:"path,attr"`
	Revision string `xml:"revision,attr"`
	Remote   string `xml:"remote,attr"`
}

type ProjectUpdate struct {
	// Time record the time when P1 was removed or P2 was added.
	// Zero value if P1/P2 are both valid(no structure changes).
	Time   time.Time
	P1, P2 *Project
}

func ParseManifestFile(file string) (*Manifest, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var m Manifest
	err = xml.Unmarshal(data, &m)
	return &m, err
}

func (m *Manifest) WriteFile(filePath string) error {
	data, err := xml.Marshal(m)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0640)
}

func GetRepoUpdates(m1, m2 *Manifest) (updates []ProjectUpdate, err error) {
	panic("implement me")
}
