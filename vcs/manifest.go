package vcs

import (
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"sort"
	"time"
)

type Manifest struct {
	XMLName  xml.Name  `xml:"manifest"`
	Remote   Remote    `xml:"remote"`
	Default  Default   `xml:"default"`
	Projects []Project `xml:"project"`
}

type Remote struct {
	Name   string `xml:"name,attr"`
	Fetch  string `xml:"fetch,attr"`
	Review string `xml:"review,attr"`
}

type Default struct {
	Remote   string `xml:"remote,attr"`
	Revision string `xml:"revision,attr"`
	SyncJ    string `xml:"sync-j,attr"`
}

type Project struct {
	Name       string     `xml:"name,attr"`
	Path       string     `xml:"path,attr,omitempty"`
	Revision   string     `xml:"revision,attr"`
	Remote     string     `xml:"remote,attr,omitempty"`
	CloneDepth string     `xml:"clone-depth,attr,omitempty"`
	LinkFile   []LinkFile `xml:"linkfile,omitempty"`
}

type LinkFile struct {
	Src  string `xml:"src,attr"`
	Dest string `xml:"dest,attr"`
}

type ProjectUpdate struct {
	// StructCTime record the time when P1 was removed or P2 was added.
	// Zero value if P1/P2 are both valid(no structure changes).
	StructCTime time.Time
	P1, P2      *Project
}

func (p *Project) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("<%s>", p.Name)
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
	data, err := xml.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	data = append([]byte(xml.Header), data...)
	return os.WriteFile(filePath, data, 0640)
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
				logrus.Infof("%v structure changes", &m1.Projects[i])
				updates = append(updates, ProjectUpdate{
					StructCTime: getTimeFn(&m1.Projects[i], &m2.Projects[j]),
					P1:          &m1.Projects[i],
					P2:          &m2.Projects[j],
				})
			} else if m1.Projects[i].Revision != m2.Projects[j].Revision {
				logrus.Infof("%v revision changes", &m1.Projects[i])
				updates = append(updates, ProjectUpdate{
					P1: &m1.Projects[i],
					P2: &m2.Projects[j],
				})
			}
		} else if m2.Projects[j].Name > m1.Projects[i].Name {
			logrus.Infof("%v removed", &m1.Projects[i])
			updates = append(updates, ProjectUpdate{
				StructCTime: getTimeFn(&m1.Projects[i], nil),
				P1:          &m1.Projects[i],
				P2:          nil,
			})
		} else {
			logrus.Infof("%v added", &m2.Projects[j])
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
	data, err := xml.MarshalIndent(m, "", "  ")
	if err != nil {
		return "", err
	}
	data = append([]byte(xml.Header), data...)
	sumByte := md5.Sum(data)
	return fmt.Sprintf("%X", sumByte), nil
}
