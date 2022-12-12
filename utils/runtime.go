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
	runtimeCache.LoadFile("gitee.cache")
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
