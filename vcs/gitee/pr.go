package gitee

import "time"

type PRSearchParam struct {
	Since time.Time
	ExtendPRSearchParam
}

type ExtendPRSearchParam struct {
	Until time.Time
}

type PRDetail struct {
}

func GetPRs(param PRSearchParam) []PRDetail {
	panic("implement me")
}
