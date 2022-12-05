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
	// StructCTime record the time when P1 was removed or P2 was added.
	// Zero value if P1/P2 are both valid(no structure changes).
	StructCTime time.Time
	P1, P2      *Project
}

func (p *Project) StructureDiff(p2 *Project) bool {
	if p == nil && p2 != nil || p != nil && p2 == nil {
		return true
	}
	return p.Name != p2.Name || p.Path != p2.Path || p.Remote != p2.Remote
}

func (p *Project) Equals(p2 *Project) bool {
	return p.Name == p2.Name && p.Path == p2.Path && p.Remote == p2.Remote && p.Revision == p2.Revision
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

func GetRepoUpdates(m1, m2 *Manifest, getTimeFn func(p1, p2 *Project) time.Time) (updates []ProjectUpdate, err error) {
	if _, err := m1.Standardize(); err != nil {
		return nil, err
	}
	if _, err := m2.Standardize(); err != nil {
		return nil, err
	}
	var j int
	for i := range m1.Projects {
		if m2.Projects[j].Name == m1.Projects[i].Name {
			if m1.Projects[i].StructureDiff(&m2.Projects[j]) {
				updates = append(updates, ProjectUpdate{
					StructCTime: getTimeFn(&m1.Projects[i], &m2.Projects[j]),
					P1:          &m1.Projects[i],
					P2:          &m2.Projects[j],
				})
			} else if m1.Projects[i].Revision != m2.Projects[j].Revision {
				updates = append(updates, ProjectUpdate{
					P1: &m1.Projects[i],
					P2: &m2.Projects[j],
				})
			}
		} else if m2.Projects[j].Name > m1.Projects[i].Name { // m1.Projects[i] has been removed in m2
			updates = append(updates, ProjectUpdate{
				StructCTime: getTimeFn(&m1.Projects[i], nil),
				P1:          &m1.Projects[i],
				P2:          nil,
			})
		} else { // m2.Projects[j] is newly added
			updates = append(updates, ProjectUpdate{
				StructCTime: getTimeFn(nil, &m1.Projects[j]),
				P1:          nil,
				P2:          &m2.Projects[j],
			})
		}
		j++
	}
	return
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
