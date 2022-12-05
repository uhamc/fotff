package vcs

import (
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"os"
	"sort"
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

func (m *Manifest) DeepCopy() (ret *Manifest) {
	ret.Default = m.Default
	ret.Projects = make([]Project, len(m.Projects))
	for i, p := range m.Projects {
		ret.Projects[i] = p
	}
	return
}

func GetRepoUpdates(m1, m2 *Manifest) (updates []ProjectUpdate, err error) {
	panic("implement me")
}

func (m *Manifest) UpdateManifestProject(name, path, remote, revision string) {
	for i, p := range m.Projects {
		if p.Name == name {
			if path != "" {
				m.Projects[i].Path = path
			}
			if remote != "" {
				m.Projects[i].Remote = remote
			}
			if revision != "" {
				m.Projects[i].Revision = revision
			}
			return
		}
	}
}

func (m *Manifest) Standardize() (string, error) {
	sort.Slice(m.Projects, func(i, j int) bool {
		return m.Projects[i].Name < m.Projects[j].Name
	})
	data, err := xml.Marshal(m)
	if err != nil {
		return "", err
	}
	sumByte := md5.Sum(data)
	return fmt.Sprintf("%X", sumByte), nil
}
