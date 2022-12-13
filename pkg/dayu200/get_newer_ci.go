package dayu200

import (
	"encoding/json"
	"fotff/utils"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type DailyBuildsQueryParam struct {
	ProjectName string `json:"projectName"`
	Branch      string `json:"branch"`
	Component   string `json:"component"`
	BuildStatus string `json:"buildStatus"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
	PageNum     int    `json:"pageNum"`
	PageSize    int    `json:"pageSize"`
}

type DailyBuildsResp struct {
	Result struct {
		DailyBuildVos []*DailyBuild `json:"dailyBuildVos"`
		Total         int           `json:"total"`
	} `json:"result"`
}

type DailyBuild struct {
	Id         string `json:"id"`
	ImgObsPath string `json:"imgObsPath"`
}

func (m *Manager) getNewerFromCI(cur string) string {
	for {
		file := func() string {
			var q = DailyBuildsQueryParam{
				ProjectName: "openharmony",
				Branch:      m.Branch,
				Component:   "dayu200",
				BuildStatus: "success",
				PageNum:     1,
				PageSize:    1,
			}
			data, err := json.Marshal(q)
			if err != nil {
				logrus.Errorf("can not marshal query param: %v", err)
				return ""
			}
			resp, err := utils.DoSimpleHttpReq(http.MethodPost, "http://ci.openharmony.cn/api/ci-backend/ci-portal/v1/dailybuilds", data, map[string]string{"Content-Type": "application/json;charset=UTF-8"})
			if err != nil {
				logrus.Errorf("can not query builds: %v", err)
				return ""
			}
			var dailyBuildsResp DailyBuildsResp
			if err := json.Unmarshal(resp, &dailyBuildsResp); err != nil {
				logrus.Errorf("can not unmarshal resp [%s]: %v", string(resp), err)
				return ""
			}
			if len(dailyBuildsResp.Result.DailyBuildVos) != 0 {
				url := dailyBuildsResp.Result.DailyBuildVos[0].ImgObsPath
				if filepath.Base(url) != cur {
					logrus.Infof("new package found, name: %s", filepath.Base(url))
					file, err := m.downloadToWorkspace(url)
					if err != nil {
						logrus.Errorf("can not download package %s: %v", url, err)
						return ""
					}
					return file
				}
			}
			return ""
		}()
		if file != "" {
			return file
		}
		time.Sleep(10 * time.Minute)
	}
}

func (m *Manager) downloadToWorkspace(url string) (string, error) {
	logrus.Infof("downloading %s", url)
	resp, err := utils.DoSimpleHttpReqRaw(http.MethodGet, url, nil, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	f, err := os.Create(filepath.Join(m.ArchiveDir, filepath.Base(url)))
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := io.CopyBuffer(f, resp.Body, make([]byte, 16*1024*1024)); err != nil {
		return "", err
	}
	logrus.Infof("%s downloaded successfully", url)
	return filepath.Base(url), nil
}
