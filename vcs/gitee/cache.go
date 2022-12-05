package gitee

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"time"
)

var respCache = cache.New(24*time.Hour, time.Hour)

func init() {
	if err := respCache.LoadFile("gitee.cache"); err != nil {
		fmt.Printf("load gitee.cache err: %v", err)
	}
}
