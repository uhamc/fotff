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

package utils

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"os"
	"path/filepath"
	"time"
)

var runtimeDir = `.fotff`

var runtimeCache = cache.New(24*time.Hour, time.Hour)

func sectionKey(section, key string) string {
	return fmt.Sprintf("__%s__%s__", section, key)
}

func init() {
	if err := os.MkdirAll(runtimeDir, 0750); err != nil {
		panic(err)
	}
	runtimeCache.LoadFile(filepath.Join(runtimeDir, "fotff.cache"))
}

func CacheGet(section string, k string) (v any, found bool) {
	return runtimeCache.Get(sectionKey(section, k))
}

func CacheSet(section string, k string, v any) error {
	runtimeCache.Set(sectionKey(section, k), v, cache.DefaultExpiration)
	return runtimeCache.SaveFile(filepath.Join(runtimeDir, "fotff.cache"))
}

func WriteRuntimeData(name string, data []byte) error {
	return os.WriteFile(filepath.Join(runtimeDir, name), data, 0640)
}

func ReadRuntimeData(name string) ([]byte, error) {
	return os.ReadFile(filepath.Join(runtimeDir, name))
}
