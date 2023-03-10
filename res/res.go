/*
 * Copyright (c) 2022 Huawei Device Co., Ltd.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package res

import (
	"fmt"
	"fotff/utils"
	"strings"
)

type Resources struct {
	DeviceSnList string `key:"device_sn_list"`
	AddrList     string `key:"build_server_addr_list" default:"127.0.0.1:22"`
	User         string `key:"build_server_user" default:"root"`
	Passwd       string `key:"build_server_password" default:"root"`
	// BuildWorkSpace must be absolute
	BuildWorkSpace string `key:"build_server_workspace" default:"/root/fotff/build_workspace"`
	devicePool     chan string
	serverPool     chan string
}

type BuildServerInfo struct {
	Addr      string
	User      string
	Passwd    string
	WorkSpace string
}

var res Resources

func init() {
	utils.ParseFromConfigFile("resources", &res)
	snList := strings.Split(res.DeviceSnList, ",")
	addrList := strings.Split(res.AddrList, ",")
	res.devicePool = make(chan string, len(snList))
	for _, sn := range snList {
		res.devicePool <- sn
	}
	res.serverPool = make(chan string, len(addrList))
	for _, addr := range addrList {
		res.serverPool <- addr
	}
}

// Fake set 'n' fake packages and build servers.
// Just for test only.
func Fake(n int) {
	var snList, addrList []string
	for i := 0; i < n; i++ {
		snList = append(snList, fmt.Sprintf("device%d", i))
		addrList = append(addrList, fmt.Sprintf("server%d", i))
	}
	res.devicePool = make(chan string, len(snList))
	for _, sn := range snList {
		res.devicePool <- sn
	}
	res.serverPool = make(chan string, len(addrList))
	for _, sn := range snList {
		res.serverPool <- sn
	}
}

func Num() int {
	if cap(res.devicePool) < cap(res.serverPool) {
		return cap(res.devicePool)
	}
	return cap(res.serverPool)
}

func DeviceList() []string {
	return strings.Split(res.DeviceSnList, ",")
}

func GetDevice() string {
	return <-res.devicePool
}

func ReleaseDevice(device string) {
	res.devicePool <- device
}

func GetBuildServer() BuildServerInfo {
	addr := <-res.serverPool
	return BuildServerInfo{
		Addr:      addr,
		User:      res.User,
		Passwd:    res.Passwd,
		WorkSpace: res.BuildWorkSpace,
	}
}

func ReleaseBuildServer(info BuildServerInfo) {
	res.serverPool <- info.Addr
}
