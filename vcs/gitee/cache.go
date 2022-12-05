package gitee

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var respCache = cache.New(24*time.Hour, time.Hour)

func init() {
	respCache.LoadFile("gitee.cache")
}
